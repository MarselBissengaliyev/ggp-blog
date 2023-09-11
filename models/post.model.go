package models

type Post struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Content     string         `json:"content"`
	PreviewUrl  string         `json:"preview_url"`
	Tags        []string       `json:"tags"`
	Comments    []Comment      `json:"comments"`
	Reactions   []PostReaction `json:"reactions"`
	UserId      string         `json:"user_id"`
	IsBanned    bool           `json:"is_banned"`
	ViewsCount  int            `json:"views_count"`
	Model
}
