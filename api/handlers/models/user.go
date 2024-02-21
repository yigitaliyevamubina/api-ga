package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int64  `json:"age"`
	Gender    int32  `json:"gender"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// User info validation
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 15), validation.Match(regexp.MustCompile("[a-z]|[A-Z][0-9]"))),
		validation.Field(&u.FirstName, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^[A-Z][a-z]*$"))),
		validation.Field(&u.LastName, validation.Required, validation.Length(5, 50), validation.Match(regexp.MustCompile("^[A-Z][a-z]*$"))),
	)
}

type UserWithPostsAndComments struct {
	Id        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Age       int64   `json:"age"`
	Gender    int32   `json:"gender"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Posts     []*Post `json:"posts"`
}

type AllUsers struct {
	Users []*UserWithPostsAndComments
}

type RegisterUserModel struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int64  `json:"age"`
	Gender    int32  `json:"gender"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
}

type VerifyUserModel struct {
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Age          int64  `json:"age"`
	Gender       int32  `json:"gender"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserModel struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Age         int64  `json:"age"`
	Gender      int32  `json:"gender"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

type LoginResp struct {
	Id          string `json:"user_id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

type RegisterRespUser struct {
	Status string `json:"message"`
}

type VerifyResponse struct {
	Status string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
