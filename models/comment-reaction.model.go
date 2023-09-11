package models

type CommentReaction struct {
	IsLiked    bool   `json:"is_liked"`
	IsDisliked bool   `json:"is_disliked"`
	CommentId  string `json:"post_id"`
	UserId     string `json:"user_id"`
	Model
}
