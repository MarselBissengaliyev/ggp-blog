package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name   string `json:"name"`
	PostId uint   `json:"post_id"`
	Post   Post   `gorm:"foreignKey:id"`
}
