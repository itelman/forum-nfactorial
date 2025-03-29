package domain

import (
	"errors"
)

var (
	ErrPostsBadRequest = errors.New("POSTS: bad request")
	ErrPostNotFound    = errors.New("DATABASE: Post not found")
)
