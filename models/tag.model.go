package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name   string `json:"name"`
	PostId uint   `json:"-"`
	Post   Post   `gorm:"foreignKey:PostId"`
}
