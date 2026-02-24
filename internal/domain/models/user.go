package models

import "time"

type User struct {
	ID                uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName         string     `json:"first_name" gorm:"size:100;not null"`
	LastName          string     `json:"last_name" gorm:"size:100;not null"`
	Email             string     `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Username          string     `json:"username" gorm:"size:50;uniqueIndex;not null"`
	PasswordHash      string     `json:"-" gorm:"size:255;not null"`
	RefreshTokenHash  string     `json:"-" gorm:"size:255"`
	OidcProvider      string     `json:"oidc_provider,omitempty" gorm:"size:50;index:idx_oidc_provider_subject"`
	OidcSubject       string     `json:"oidc_subject,omitempty" gorm:"size:255;index:idx_oidc_provider_subject"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	IsActive          bool       `json:"is_active" gorm:"not null;default:true"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy         uint64     `json:"created_by" gorm:"index"`

	// Relationships
	TaskGroups []*TaskGroup `json:"task_groups,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (User) TableName() string {
	return "users"
}
