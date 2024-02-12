package models

type Comment struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	PostId  string `json:"post_id"`
	OwnerId string `json:"owner_id"`
}

type Comments struct {
	Comments []*Comment
}
