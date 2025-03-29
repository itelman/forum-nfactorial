package domain

import (
	"errors"
)

var (
	ErrPostReactionsBadRequest = errors.New("POST REACTIONS: bad request")
	ErrPostReactionNotFound    = errors.New("DATABASE: Post reaction not found")
)
