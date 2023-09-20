package models

type Report struct {
	Reason string `json:"reason"`
	UserId uint   `json:"user_id"`
	PostId uint   `json:"post_id"`
	Model
}
