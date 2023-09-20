package models

type Tag struct {
	PostId uint   `json:"post_id"`
	Name   string `json:"name"`
	Model
}
