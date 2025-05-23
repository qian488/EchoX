package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:50;not null;unique"`
	Color     string         `json:"color" gorm:"size:7;default:'#409EFF'"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Records   []Record       `json:"records" gorm:"many2many:record_tags;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
