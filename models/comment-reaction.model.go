package models

type CommentReaction struct {
	IsLiked    bool `json:"is_liked"`
	IsDisliked bool `json:"is_disliked"`
	CommentId  uint `json:"post_id"`
	UserId     uint `json:"user_id"`
	Model
}
