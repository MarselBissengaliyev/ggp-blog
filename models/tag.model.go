package models

type Tag struct {
	Name   string `json:"name"`
	PostId uint   `json:"post_id"`
	Model
}
