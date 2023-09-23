package models

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
