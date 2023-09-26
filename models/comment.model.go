package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserId  uint   `json:"user_id"`
	User    User   `gorm:"primaryKey:id"`
	PostId  uint   `json:"post_id"`
	Post    Post   `gorm:"primary_key:id"`
}
