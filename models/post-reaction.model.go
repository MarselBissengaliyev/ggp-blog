package models

import "gorm.io/gorm"

type PostReaction struct {
	gorm.Model
	IsLiked    bool `json:"is_liked"`
	IsDisliked bool `json:"is_disliked"`
	PostId     uint `json:"post_id"`
	Post       Post `gorm:"foreignKey:PostId"`
	UserId     uint `json:"user_id"`
	User       User `gorm:"foreignKey:UserId"`
}
