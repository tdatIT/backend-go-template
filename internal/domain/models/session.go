package models

import "time"

type Session struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint64     `json:"user_id" gorm:"not null;index"`
	RefreshJTI string     `json:"refresh_jti" gorm:"size:255;not null;uniqueIndex"`
	UserAgent  string     `json:"user_agent,omitempty" gorm:"size:255"`
	IPAddress  string     `json:"ip_address,omitempty" gorm:"size:50"`
	IsActive   bool       `json:"is_active" gorm:"not null;default:true"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
