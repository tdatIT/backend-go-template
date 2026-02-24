package models

import "time"

type Session struct {
	ID         string    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint64    `json:"user_id" gorm:"not null;index"`
	RefreshJTI string    `json:"refresh_jti" gorm:"size:255;not null;uniqueIndex"`
	UserAgent  string    `json:"user_agent,omitempty" gorm:"size:255"`
	DeviceID   string    `json:"device_id,omitempty" gorm:"size:255"`
	Platform   string    `json:"platform,omitempty" gorm:"size:100"`
	IP         string    `json:"ip,omitempty" gorm:"size:50"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
