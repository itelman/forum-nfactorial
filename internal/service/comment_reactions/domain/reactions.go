package domain

import (
	"errors"
)

var (
	ErrCommentReactionsBadRequest = errors.New("COMMENT REACTIONS: bad request")
	ErrCommentReactionNotFound    = errors.New("DATABASE: Post reaction not found")
)
