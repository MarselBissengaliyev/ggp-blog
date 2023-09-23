package models

type Comment struct {
	Content string `json:"content"`
	UserId  uint   `json:"user_id"`
	PostId  uint   `json:"post_id"`
	Model
}