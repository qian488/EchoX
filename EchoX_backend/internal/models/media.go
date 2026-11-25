package models

import (
	"time"

	"gorm.io/gorm"
)

type MediaFile struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	Filename     string         `gorm:"size:255;not null" json:"filename"`
	OriginalName string         `gorm:"size:255;not null" json:"original_name"`
	MimeType     string         `gorm:"size:100" json:"mime_type"`
	Size         int64          `gorm:"not null" json:"size"`
	Type         string         `gorm:"size:20;not null" json:"type"`
	URL          string         `gorm:"size:512;not null" json:"url"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
