package models

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Title     string         `gorm:"size:100;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Type      string         `gorm:"size:20;not null" json:"type"`
	MediaURL  string         `gorm:"size:255" json:"media_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Tags      []Tag          `json:"tags" gorm:"many2many:record_tags;"`
}
