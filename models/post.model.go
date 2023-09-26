package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title         string         `json:"title"`
	Slug          string         `json:"slug"`
	Description   string         `json:"description"`
	Content       string         `json:"content"`
	PreviewUrl    string         `json:"preview_url"`
	UserId        uint           `json:"user_id"`
	User          User           `gorm:"foreignKey:UserId"`
	IsBanned      bool           `json:"is_banned"`
	ViewsCount    int            `json:"views_count"`
	PostReactions []PostReaction `gorm:"foreignKey:post_id"`
}