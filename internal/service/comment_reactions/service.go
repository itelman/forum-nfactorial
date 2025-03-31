package comment_reactions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/pkg/requests"
	"net/http"
)

var (
	ErrCommentNotFound                  = errors.New("DATABASE: Comment not found")
	ErrCommentReactionsBadRequest       = errors.New("COMMENT REACTIONS: bad request")
	ErrCommentReactionsUserUnauthorized = errors.New("COMMENT REACTIONS: user unauthorized")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("COMMENT REACTIONS (API): %s", status)
}

type Service interface {
	CreateCommentReaction(ctx context.Context, input *CreateCommentReactionInput) error
}

type service struct {
	postsEndpoint string
}

func NewService(opts ...Option) *service {
	svc := &service{}
	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

type Option func(*service)

func WithAPI(apiLink string) Option {
	return func(s *service) {
		s.postsEndpoint = apiLink + "/posts"
	}
}

type CreateCommentReactionInput struct {
	CommentID int
	PostID    int
	IsLike    int
}

func (s *service) CreateCommentReaction(ctx context.Context, input *CreateCommentReactionInput) error {
	reqInput := struct {
		IsLike int `json:"is_like"`
	}{input.IsLike}

	reqBody, err := json.Marshal(reqInput)
	if err != nil {
		return err
	}

	resp, err := requests.SendRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%d/comments/%d/react", s.postsEndpoint, input.PostID, input.CommentID),
		bytes.NewBuffer(reqBody),
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)),
		},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		return ErrCommentReactionsBadRequest
	} else if resp.StatusCode == http.StatusUnauthorized {
		return ErrCommentReactionsUserUnauthorized
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return ErrAPIUnhandled(resp.Status)
	}

	return nil
}
