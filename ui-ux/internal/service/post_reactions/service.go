package post_reactions

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
	ErrPostNotFound                  = errors.New("DATABASE: Post not found")
	ErrPostReactionsBadRequest       = errors.New("POST REACTIONS: bad request")
	ErrPostReactionsUserUnauthorized = errors.New("POST REACTIONS: user unauthorized")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("POST REACTIONS (API): %s", status)
}

type Service interface {
	CreatePostReaction(ctx context.Context, input *CreatePostReactionInput) error
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

type CreatePostReactionInput struct {
	PostID int
	IsLike int
}

func (s *service) CreatePostReaction(ctx context.Context, input *CreatePostReactionInput) error {
	reqInput := struct {
		IsLike int `json:"is_like"`
	}{input.IsLike}

	reqBody, err := json.Marshal(reqInput)
	if err != nil {
		return err
	}

	resp, err := requests.SendRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%d/react", s.postsEndpoint, input.PostID),
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
		return ErrPostReactionsBadRequest
	} else if resp.StatusCode == http.StatusUnauthorized {
		return ErrPostReactionsUserUnauthorized
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return ErrAPIUnhandled(resp.Status)
	}

	return nil
}
