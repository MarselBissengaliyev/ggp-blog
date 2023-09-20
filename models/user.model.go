package models

import (
	"database/sql"
	"time"
)

// Updated password
type User struct {
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	UserName       string         `json:"user_name" gorm:"unique"`
	Role           string         `json:"role"`
	Email          string         `json:"email" gorm:"unique"`
	EmailCode      sql.NullString `json:"email_code" gorm:"unique"`
	EmailConfirmed bool           `json:"email_confirmed"`
	IsBanned       bool           `json:"is_banned"`
	BanTime        time.Time      `json:"ban_time"`
	Password       string         `json:"password"`
	Model
}
