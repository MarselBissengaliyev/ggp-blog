package models

import "gorm.io/gorm"

type CommentReaction struct {
	gorm.Model
	IsLiked    bool `json:"is_liked"`
	IsDisliked bool `json:"is_disliked"`
	PostId     uint `json:"post_id"`
	Post       Post `gorm:"primary_key:id"`
	UserId     uint `json:"user_id"`
	User       User `gorm:"primaryKey:id"`
}
