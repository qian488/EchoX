package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"echox/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PasswordController struct {
	DB *gorm.DB
}

func NewPasswordController(db *gorm.DB) *PasswordController {
	return &PasswordController{DB: db}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// generateResetToken 使用 crypto/rand 生成 32 字节随机数，编码为 64 位 hex 字符串
func generateResetToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// ForgotPassword 忘记密码：请求重置链接
func (pc *PasswordController) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 统一成功响应，防止邮箱枚举：无论用户是否存在都返回相同 message
	response := gin.H{"message": "如果该邮箱已注册，重置链接已发送"}

	var user models.User
	if err := pc.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// 用户不存在：不暴露任何信息，直接返回统一响应
		c.JSON(http.StatusOK, response)
		return
	}

	// 生成新 token 前，先将该用户之前未使用的 token 全部标记为已用（失效）
	if err := pc.DB.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used = ?", user.ID, false).
		Update("used", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "使旧令牌失效失败"})
		return
	}

	token, err := generateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成重置令牌失败"})
		return
	}

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		Used:      false,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	if err := pc.DB.Create(&resetToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存重置令牌失败"})
		return
	}

	// ===== 开发模式 =====
	// 本项目是本地比赛作品，未配置邮件服务，因此把 reset_token 一并返回，
	// 方便前端展示给用户测试。生产环境必须改为：通过邮件发送含 token 的重置链接，
	// 并且绝不在此处返回 token，否则任何人都可重置任意账户密码。
	response["reset_token"] = token

	c.JSON(http.StatusOK, response)
}

// ResetPassword 重置密码：校验令牌并更新密码
func (pc *PasswordController) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 令牌必须存在、未使用、未过期
	var resetToken models.PasswordResetToken
	if err := pc.DB.Where("token = ? AND used = ? AND expires_at > ?", req.Token, false, time.Now()).
		First(&resetToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "重置令牌无效或已过期"})
		return
	}

	// 关联用户必须存在（First 会自动排除软删除用户）
	var user models.User
	if err := pc.DB.First(&user, resetToken.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "重置令牌无效或已过期"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加密密码"})
		return
	}

	// 用事务保证「更新密码」与「令牌标记已用」同步成功，避免令牌被重复使用
	err = pc.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).
			Update("password", string(hashedPassword)).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.PasswordResetToken{}).Where("id = ?", resetToken.ID).
			Update("used", true).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重置密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}
