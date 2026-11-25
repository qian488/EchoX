package controllers

import (
	"net/http"
	"os"
	"time"

	"echox/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名或邮箱是否已存在
	var existingUser models.User
	if uc.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加密密码"})
		return
	}
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// 创建用户
	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
}

// Login 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var loginForm LoginRequest
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if uc.DB.Where("username = ?", loginForm.Username).First(&user).Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token有效期24小时
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT Secret未配置"})
		return
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成Token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// GetProfile 获取个人信息
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var user models.User
	// 查找用户，排除密码字段
	if uc.DB.Select("id, username, email, created_at, updated_at").First(&user, userID).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取用户的记录统计信息
	var totalRecords int64
	uc.DB.Model(&models.Record{}).Where("user_id = ?", userID).Count(&totalRecords)

	var textRecords int64
	uc.DB.Model(&models.Record{}).Where("user_id = ? AND type = ?", userID, "text").Count(&textRecords)

	var imageRecords int64
	uc.DB.Model(&models.Record{}).Where("user_id = ? AND type = ?", userID, "image").Count(&imageRecords)

	var videoRecords int64
	uc.DB.Model(&models.Record{}).Where("user_id = ? AND type = ?", userID, "video").Count(&videoRecords)

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"stats": gin.H{
			"total": totalRecords,
			"text":  textRecords,
			"image": imageRecords,
			"video": videoRecords,
		},
	})
}

// UpdateProfile 更新个人信息
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var user models.User
	if uc.DB.First(&user, userID).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var updateForm struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
	}

	if err := c.ShouldBindJSON(&updateForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名或邮箱是否已存在（排除当前用户）
	if updateForm.Username != nil || updateForm.Email != nil {
		var existingUser models.User
		query := uc.DB.Where("id != ?", userID)
		if updateForm.Username != nil && updateForm.Email != nil {
			query = query.Where("username = ? OR email = ?", *updateForm.Username, *updateForm.Email)
		} else if updateForm.Username != nil {
			query = query.Where("username = ?", *updateForm.Username)
		} else if updateForm.Email != nil {
			query = query.Where("email = ?", *updateForm.Email)
		}

		if query.First(&existingUser).Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱已存在"})
			return
		}
	}

	// 更新字段
	if updateForm.Username != nil {
		user.Username = *updateForm.Username
	}
	if updateForm.Email != nil {
		user.Email = *updateForm.Email
	}

	// 保存更新
	if uc.DB.Save(&user).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}
