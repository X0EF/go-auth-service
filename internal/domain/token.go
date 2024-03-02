package domain

import (
	"time"

	"gorm.io/gorm"
)

type TokenType string

const (
	ACCESS         TokenType = "access"
	REFRESH        TokenType = "refresh"
	RESET_PASSWORD TokenType = "reset_password"
	VERIFY_EMAIL   TokenType = "verify_email"
	VERIFY_PHONE   TokenType = "verify_phone"
)

type Token struct {
	ID        uint `gorm:"primaryKey;autoIncrement:true;not null" json:"id"`
	UserID    uint `                                              json:"user_id"`
	User      User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Type      TokenType
	Hash      string         `gorm:"unique"`
	ExpiresAt time.Time      `                                              `
	CreatedAt time.Time      `                                              json:"created_at"`
	UpdatedAt time.Time      `                                              json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                  json:"-"`
}
