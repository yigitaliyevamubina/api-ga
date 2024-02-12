package models

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int64  `json:"age"`
	Gender    int32  `json:"gender"`
}

type UserWithPostsAndComments struct {
	Id        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Age       int64   `json:"age"`
	Gender    int32   `json:"gender"`
	Posts     []*Post `json:"posts"`
}

type AllUsers struct {
	Users []*UserWithPostsAndComments
}
