package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	UserName       string     `json:"user_name" gorm:"unique"`
	Role           string     `json:"-"`
	Email          string     `json:"email" gorm:"unique"`
	EmailCode      *string    `json:"-" gorm:"unique"`
	EmailConfirmed bool       `json:"email_confirmed"`
	IsBanned       bool       `json:"is_banned"`
	BanTime        *time.Time `json:"ban_time"`
	Password       string     `json:"password"`
}
