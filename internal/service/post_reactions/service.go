package post_reactions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"net/http"
)

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

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%d/react", s.postsEndpoint, input.PostID),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)))

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if apiResp.StatusCode != http.StatusOK {
		return errors.New("POST REACTIONS: /CREATE - API ERROR")
	}

	return nil
}
