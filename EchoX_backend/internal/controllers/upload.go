package controllers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"echox/internal/config"
	"echox/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadController struct {
	DB *gorm.DB
}

func NewUploadController(db *gorm.DB) *UploadController {
	return &UploadController{DB: db}
}

// MediaResponse 是上传与查询接口共用的响应结构。
type MediaResponse struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

func toMediaResponse(m models.MediaFile) MediaResponse {
	return MediaResponse{
		ID:       m.ID,
		URL:      m.URL,
		Filename: m.OriginalName,
		MimeType: m.MimeType,
		Size:     m.Size,
		Type:     m.Type,
	}
}

// isBodyTooLarge 判断 multipart 解析错误是否由请求体超限引起。
func isBodyTooLarge(err error) bool {
	var maxErr *http.MaxBytesError
	return errors.As(err, &maxErr) || strings.Contains(err.Error(), "request body too large")
}

// Upload 处理 multipart/form-data 文件上传：字段 file 必填，type 可选（text|image|video，缺省 image）。
func (uc *UploadController) Upload(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或登录已失效"})
		return
	}

	// 限制整个请求体，给 multipart 边界与表单字段预留 1MB 余量。
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxUploadSize+(1<<20))

	fileType := c.PostForm("type")
	if fileType == "" {
		fileType = "image"
	}
	if fileType != "text" && fileType != "image" && fileType != "video" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 type，仅支持 text/image/video"})
		return
	}

	header, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少文件字段 file"})
			return
		}
		if isBodyTooLarge(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过 100MB 限制"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析上传文件失败"})
		return
	}

	if header.Size > config.MaxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过 100MB 限制"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !config.IsAllowedExtension(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，仅支持 png/jpg/jpeg/gif/webp/mp4/webm/mov/avi"})
		return
	}

	originalName := filepath.Base(header.Filename)
	storedName := config.GenerateFilename(header.Filename)
	mimeType := config.GetMimeType(ext)
	url := "/uploads/" + storedName

	if err := config.EnsureUploadDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	src, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取上传文件失败"})
		return
	}
	defer src.Close()

	// O_EXCL 保证随机名即使碰撞也不会覆盖已有文件。
	destPath := filepath.Join(config.GetUploadDir(), storedName)
	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目标文件失败"})
		return
	}

	_, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil || closeErr != nil {
		os.Remove(destPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入文件失败"})
		return
	}

	media := models.MediaFile{
		UserID:       uid,
		Filename:     storedName,
		OriginalName: originalName,
		MimeType:     mimeType,
		Size:         header.Size,
		Type:         fileType,
		URL:          url,
	}
	if err := uc.DB.Create(&media).Error; err != nil {
		os.Remove(destPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存媒体记录失败"})
		return
	}

	c.JSON(http.StatusCreated, toMediaResponse(media))
}

// GetMedia 返回媒体信息，仅文件所有者可访问。
func (uc *UploadController) GetMedia(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的媒体ID"})
		return
	}

	var media models.MediaFile
	if err := uc.DB.First(&media, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "媒体文件不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取媒体信息失败"})
		return
	}

	userID, _ := c.Get("user_id")
	if media.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此文件"})
		return
	}

	c.JSON(http.StatusOK, toMediaResponse(media))
}

// DeleteMedia 校验归属后删除数据库记录（软删）与物理文件。
func (uc *UploadController) DeleteMedia(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的媒体ID"})
		return
	}

	var media models.MediaFile
	if err := uc.DB.First(&media, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "媒体文件不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取媒体信息失败"})
		return
	}

	userID, _ := c.Get("user_id")
	if media.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此文件"})
		return
	}

	if err := os.Remove(filepath.Join(config.GetUploadDir(), media.Filename)); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除物理文件失败"})
		return
	}

	if err := uc.DB.Delete(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除媒体记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件已删除"})
}
