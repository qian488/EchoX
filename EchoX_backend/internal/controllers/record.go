package controllers

import (
	"net/http"
	"strconv"

	"echox/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecordController struct {
	DB *gorm.DB
}

func NewRecordController(db *gorm.DB) *RecordController {
	return &RecordController{DB: db}
}

type CreateRecordRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	Type     string `json:"type" binding:"required,oneof=text image video"`
	MediaURL string `json:"media_url"`
	TagIDs   []uint `json:"tag_ids"`
}

type UpdateRecordDto struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Type     string `json:"type" binding:"omitempty,oneof=text image video"`
	MediaURL string `json:"media_url"`
	TagIDs   *[]uint `json:"tag_ids"`
}

// CreateRecord 创建新记录
func (rc *RecordController) CreateRecord(c *gin.Context) {
	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置用户ID
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	record := models.Record{
		Title:    req.Title,
		Content:  req.Content,
		Type:     req.Type,
		MediaURL: req.MediaURL,
		UserID:   uid,
	}

	// 校验标签归属，防止把别人的标签挂到自己的记录上
	var tags []models.Tag
	if len(req.TagIDs) > 0 {
		if err := rc.DB.Where("id IN ? AND user_id = ?", req.TagIDs, uid).Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询标签失败"})
			return
		}
		if len(tags) != len(req.TagIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "部分标签不存在或无权使用"})
			return
		}
		record.Tags = tags
	}

	if err := rc.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建记录失败"})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// GetRecords 获取用户的所有记录
func (rc *RecordController) GetRecords(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var records []models.Record

	if err := rc.DB.Where("user_id = ?", userID).Preload("Tags").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// GetRecord 获取单条记录
func (rc *RecordController) GetRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	var record models.Record
	if err := rc.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	// 验证所有权
	userID, _ := c.Get("user_id")
	if record.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此记录"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// UpdateRecord 更新记录
func (rc *RecordController) UpdateRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	var record models.Record
	if err := rc.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	// 验证所有权
	userID, _ := c.Get("user_id")
	if record.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此记录"})
		return
	}

	var updateData UpdateRecordDto
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if updateData.Title != "" {
		record.Title = updateData.Title
	}
	if updateData.Content != "" {
		record.Content = updateData.Content
	}
	if updateData.Type != "" {
		record.Type = updateData.Type
	}
	if updateData.MediaURL != "" {
		record.MediaURL = updateData.MediaURL
	}

	if updateData.TagIDs != nil {
		var tags []models.Tag
		if len(*updateData.TagIDs) > 0 {
			if err := rc.DB.Where("id IN ? AND user_id = ?", *updateData.TagIDs, userID).Find(&tags).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询标签失败"})
				return
			}
			if len(tags) != len(*updateData.TagIDs) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "部分标签不存在或无权使用"})
				return
			}
		}
		if err := rc.DB.Model(&record).Association("Tags").Replace(tags); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
			return
		}
	}

	if err := rc.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新记录失败"})
		return
	}

	rc.DB.Preload("Tags").First(&record, record.ID)
	c.JSON(http.StatusOK, record)
}

// DeleteRecord 删除记录
func (rc *RecordController) DeleteRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	var record models.Record
	if err := rc.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	// 验证所有权
	userID, _ := c.Get("user_id")
	if record.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此记录"})
		return
	}

	if err := rc.DB.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "记录已删除"})
}
