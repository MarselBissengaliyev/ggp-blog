package models

type PostReaction struct {
	IsLiked    bool `json:"is_liked"`
	IsDisliked bool `json:"is_disliked"`
	PostId     uint `json:"post_id"`
	UserId     uint `json:"user_id"`
	Model
}