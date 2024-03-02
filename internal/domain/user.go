package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey;autoIncrement:true;not null" json:"id"`
	Username        string         `gorm:"unique;not null"                        json:"username"`
	Email           string         `gorm:"unique;not null"                        json:"email"`
	Password        string         `gorm:"not null"                               json:"-"`
	Phone           *string        `gorm:"unique;"                                json:"phone"`
	Firstname       string         `                                              json:"first_name"`
	Lastname        string         `                                              json:"last_name"`
	AvatarURL       string         `                                              json:"avatar_url"`
	Role            Role           `                                              json:"role"`
	IsEmailVerified bool           `                                              json:"is_email_verified"`
	LastLoggedInAt  time.Time      `                                              json:"last_logged_in_at"`
	VerifiedAt      time.Time      `                                              json:"verified_at"`
	CreatedAt       time.Time      `                                              json:"created_at"`
	UpdatedAt       time.Time      `                                              json:"updated_at,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index"                                  json:"-"`
}
