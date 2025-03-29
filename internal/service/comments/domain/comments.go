package domain

import (
	"errors"
)

var (
	ErrCommentsBadRequest = errors.New("COMMENTS: bad request")
	ErrCommentNotFound    = errors.New("DATABASE: Comment not found")
)
