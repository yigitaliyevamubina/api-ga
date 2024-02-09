package models

type CommentLike struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	CommentId string `json:"comment_id"`
}

type PostLike struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
}
