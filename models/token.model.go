package models

import "gorm.io/gorm"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	UserId       string `json:"user_id"`
	gorm.Model
}
