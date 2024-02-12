package models

type Post struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url"`
	OwnerId  string `json:"owner_id"`
}

type RespPost struct {
	Id       string     `json:"id"`
	Title    string     `json:"title"`
	ImageUrl string     `json:"image_url"`
	OwnerId  string     `json:"owner_id"`
	Comments []*Comment `json:"comments"`
}

type Posts struct {
	Posts []*Post `json:"posts"`
}
