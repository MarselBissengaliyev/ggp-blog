package models

type Comment struct {
	Content   string            `json:"content"`
	UserId    string            `json:"user_id"`
	PostId    string            `json:"post_id"`
	Reactions []CommentReaction `json:"reactions"`
	Model
}
