package controllers

import (
	"net/http"
	"time"

	"echox/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DiscoveryController struct {
	DB *gorm.DB
}

func NewDiscoveryController(db *gorm.DB) *DiscoveryController {
	return &DiscoveryController{DB: db}
}

const dateLayout = "2006-01-02"

// Search 按当前用户搜索记录，支持关键词、类型、标签和日期范围过滤
func (dc *DiscoveryController) Search(c *gin.Context) {
	userID, _ := c.Get("user_id")

	query := dc.DB.Model(&models.Record{}).Where("records.user_id = ?", userID)

	// 关键词：标题或内容模糊匹配
	if q := c.Query("q"); q != "" {
		pattern := "%" + q + "%"
		query = query.Where("(records.title LIKE ? OR records.content LIKE ?)", pattern, pattern)
	}

	// 类型：精确过滤
	if recordType := c.Query("type"); recordType != "" {
		switch recordType {
		case "text", "image", "video":
			query = query.Where("records.type = ?", recordType)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "type 参数无效"})
			return
		}
	}

	// 标签：通过 record_tags 与 tags 关联精确匹配标签名
	if tagName := c.Query("tag"); tagName != "" {
		query = query.
			Joins("JOIN record_tags ON record_tags.record_id = records.id").
			Joins("JOIN tags ON tags.id = record_tags.tag_id").
			Where("tags.name = ? AND tags.deleted_at IS NULL", tagName).
			Distinct()
	}

	// 日期范围：对 created_at 过滤
	if start := c.Query("start"); start != "" {
		t, err := time.Parse(dateLayout, start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start 日期格式无效"})
			return
		}
		query = query.Where("records.created_at >= ?", t)
	}
	if end := c.Query("end"); end != "" {
		t, err := time.Parse(dateLayout, end)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end 日期格式无效"})
			return
		}
		query = query.Where("records.created_at <= ?", t.AddDate(0, 0, 1).Add(-time.Second))
	}

	records := make([]models.Record, 0)
	if err := query.Order("records.created_at DESC").Preload("Tags").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索记录失败"})
		return
	}

	c.JSON(http.StatusOK, records)
}

type timelineGroup struct {
	Date    string          `json:"date"`
	Count   int             `json:"count"`
	Records []models.Record `json:"records"`
}

// Timeline 聚合当前用户的记录为时间轴，按日期倒序返回
func (dc *DiscoveryController) Timeline(c *gin.Context) {
	userID, _ := c.Get("user_id")

	group := c.Query("group")
	if group == "" {
		group = "day"
	}
	var layout string
	switch group {
	case "day":
		layout = "2006-01-02"
	case "month":
		layout = "2006-01"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "group 参数无效"})
		return
	}

	records := make([]models.Record, 0)
	if err := dc.DB.Model(&models.Record{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("Tags").
		Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取时间轴失败"})
		return
	}

	groups := make(map[string]*timelineGroup)
	order := make([]string, 0)
	for _, record := range records {
		key := record.CreatedAt.Format(layout)
		g, ok := groups[key]
		if !ok {
			g = &timelineGroup{Date: key, Records: make([]models.Record, 0)}
			groups[key] = g
			order = append(order, key)
		}
		g.Records = append(g.Records, record)
		g.Count++
	}

	result := make([]timelineGroup, 0, len(order))
	for _, key := range order {
		result = append(result, *groups[key])
	}

	c.JSON(http.StatusOK, result)
}
