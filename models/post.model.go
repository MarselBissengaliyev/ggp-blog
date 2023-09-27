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
	ViewsCount    int            `json:"views_count"`
	PostReactions []PostReaction `json:"-" gorm:"foreignKey:post_id"`
	Comments      []Comment      `json:"-" gorm:"foreignKey:post_id"`
	Tags          []Tag          `json:"tags" gorm:"foreignKey:post_id"`
}
