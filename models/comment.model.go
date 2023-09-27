package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserId  uint   `json:"user_id"`
	User    User   `gorm:"foreignKey:UserId"`
	PostId  uint   `json:"post_id"`
	Post    Post   `gorm:"foreignKey:PostId"`
}
