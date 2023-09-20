package models

import (
	"fmt"

	"gorm.io/gorm"
)

type PostReaction struct {
	IsLiked    bool `json:"is_liked"`
	IsDisliked bool `json:"is_disliked"`
	PostId     uint `json:"post_id"`
	UserId     uint `json:"user_id"`
	Model
}

func (pr *PostReaction) Author(db *gorm.DB) string {
	var author User

	err := db.First(&author, fmt.Sprintf("id = %d", pr.UserId)).Error
	if err != nil {
		return ""
	}

	return author.UserName
}

func (pr *PostReaction) PostSlug(db *gorm.DB) string {
	var post Post

	err := db.First(&post, fmt.Sprintf("id = %d", pr.PostId)).Error
	if err != nil {
		return ""
	}

	return post.Slug
}
