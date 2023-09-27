package models

import "gorm.io/gorm"

type Report struct {
	gorm.Model
	Reason string `json:"reason"`
	PostId uint   `json:"post_id"`
	Post   Post   `gorm:"foreignKey:PostId"`
	UserId uint   `json:"user_id"`
	User   User   `gorm:"foreignKey:UserId"`
}
