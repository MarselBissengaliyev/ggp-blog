package models

import "time"

type User struct {
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	UserName       string    `json:"user_name"`
	Role           string    `json:"role"`
	Email          string    `json:"email"`
	EmailCode      string    `json:"email_code"`
	EmailConfirmed bool      `json:"email_confirmed"`
	IsBanned       bool      `json:"is_banned"`
	BanTime        time.Time `json:"ban_time"`
	Model
}
