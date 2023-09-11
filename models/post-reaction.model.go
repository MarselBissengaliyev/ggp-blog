package models

type PostReaction struct {
	IsLiked    bool   `json:"is_liked"`
	IsDisliked bool   `json:"is_disliked"`
	PostId     string `json:"post_id"`
	UserId     string `json:"user_id"`
	Model
}
