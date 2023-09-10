package models

import "gorm.io/gorm"

type User struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	UserName       string `json:"user_name"`
	Role           string `json:"role"`
	Email          string `json:"email"`
	EmailCode      string
	EmailConfirmed bool
	gorm.Model
}
