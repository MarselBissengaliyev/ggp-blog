package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Post struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Content     string `json:"content"`
	PreviewUrl  string `json:"preview_url"`
	UserId      uint   `json:"user_id"`
	IsBanned    bool   `json:"is_banned"`
	ViewsCount  int    `json:"views_count"`
	Model
}

func (p *Post) ReactionsCount(db *gorm.DB, postId uint) int64 {
	var count int64

	err := db.Model(PostReaction{}).Where("post_id = ?", p.ID).Count(&count).Error

	if err != nil {
		return 0
	}

	return count
}

func (p *Post) Author(db *gorm.DB) string {
	var author User

	err := db.First(&author, fmt.Sprintf("id = %d", p.UserId)).Error
	if err != nil {
		return ""
	}

	return author.UserName
}
