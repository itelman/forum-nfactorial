package domain

import (
	"database/sql"
	"errors"

	"github.com/itelman/forum/internal/dto"
)

type PostsRepository interface {
	Get(input GetPostInput) (*dto.Post, error)
	UpdateReactionsCount(tx *sql.Tx, input UpdatePostReactionsCountInput) error
}

type UpdatePostReactionsCountInput struct {
	PostID int
}

type GetPostInput struct {
	ID int
}

var (
	ErrPostsBadRequest = errors.New("POSTS: bad request")
	ErrPostNotFound    = errors.New("DATABASE: Post not found")
)
