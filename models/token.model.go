package models

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	UserId       uint   `json:"user_id"`
	User         User   `gorm:"foreignKey:UserId"`
}
