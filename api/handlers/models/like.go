package models

type CommentLike struct {
	UserId    string `json:"user_id"`
	CommentId string `json:"comment_id"`
}

type PostLike struct {
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
}

type Status struct {
	Liked bool `json:"liked"`
}

type ResponseLikePost struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	ImageUrl string  `json:"image_url"`
	OwnerId  string  `json:"owner_id"`
	Likes    []*User `json:"likes"`
}

type ResponseLikeComment struct {
	Id        string  `json:"id"`
	Content   string  `json:"content"`
	PostId    string  `json:"post_id"`
	OwnerId   string  `json:"owner_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Likes     []*User `json:"likes"`
}
