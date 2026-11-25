package models

import "time"

// PasswordResetToken 密码重置令牌
// 每个令牌对应一次重置请求，30 分钟内有效，使用后标记失效。
type PasswordResetToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"size:64;not null;unique;index" json:"token"`
	Used      bool      `gorm:"default:false" json:"used"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
