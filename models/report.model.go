package models

type Report struct {
	Reason string `json:"reason"`
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
	Model
}
