package models

type Tag struct {
	PostId string   `json:"post_id"`
	Name   string   `json:"name"`
	Posts  []string `json:"posts"`
	Model
}
