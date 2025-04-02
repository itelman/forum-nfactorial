package users

import (
	"net/http"
)

func DecodeSignupUser(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrUsersBadRequest
	}

	return &SignupUserInput{
		Username: r.PostForm.Get("username"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}, nil
}

func DecodeLoginUser(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrUsersBadRequest
	}

	return &LoginUserInput{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	}, nil
}
