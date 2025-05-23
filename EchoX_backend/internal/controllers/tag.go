package controllers

import (
	"net/http"
	"strconv"

	"echox/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TagController struct {
	DB *gorm.DB
}

func NewTagController(db *gorm.DB) *TagController {
	return &TagController{DB: db}
}

// CreateTag 创建新标签
func (tc *TagController) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置用户ID
	userID, _ := c.Get("user_id")
	tag.UserID = userID.(uint)

	// 检查标签名是否已存在
	var existingTag models.Tag
	if err := tc.DB.Where("name = ? AND user_id = ?", tag.Name, tag.UserID).First(&existingTag).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签名已存在"})
		return
	}

	if err := tc.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTags 获取用户的所有标签
func (tc *TagController) GetTags(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var tags []models.Tag

	if err := tc.DB.Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// UpdateTag 更新标签
func (tc *TagController) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var tag models.Tag
	if err := tc.DB.First(&tag, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	// 验证所有权
	userID, _ := c.Get("user_id")
	if tag.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此标签"})
		return
	}

	var updateData struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查新标签名是否与其他标签重复
	if updateData.Name != tag.Name {
		var existingTag models.Tag
		if err := tc.DB.Where("name = ? AND user_id = ? AND id != ?", updateData.Name, userID, id).First(&existingTag).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "标签名已存在"})
			return
		}
	}

	tag.Name = updateData.Name
	tag.Color = updateData.Color

	if err := tc.DB.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag 删除标签
func (tc *TagController) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var tag models.Tag
	if err := tc.DB.First(&tag, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	// 验证所有权
	userID, _ := c.Get("user_id")
	if tag.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此标签"})
		return
	}

	if err := tc.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标签已删除"})
}
