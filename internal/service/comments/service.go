package comments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"io/ioutil"
	"net/http"
)

type Service interface {
	CreateComment(ctx context.Context, input *CreateCommentInput) error
	GetComment(input *GetCommentInput) (*GetCommentResponse, error)
	UpdateComment(ctx context.Context, input *UpdateCommentInput) error
	DeleteComment(ctx context.Context, input *DeleteCommentInput) error
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

type CreateCommentInput struct {
	PostID  int
	Content string
}

func (s *service) CreateComment(ctx context.Context, input *CreateCommentInput) error {
	reqInput := struct {
		Content string `json:"content"`
	}{input.Content}

	reqBody, err := json.Marshal(reqInput)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%d/comments", s.postsEndpoint, input.PostID),
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
		return errors.New("COMMENTS: /CREATE - API ERROR")
	}

	return nil
}

type GetCommentInput struct {
	PostID int
	ID     int
}

type GetCommentResponse struct {
	Comment *dto.Comment
}

func (s *service) GetComment(input *GetCommentInput) (*GetCommentResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%d/comments/%d", s.postsEndpoint, input.PostID, input.ID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if apiResp.StatusCode != http.StatusOK {
		return nil, errors.New("COMMENTS: /GET - API ERROR")
	}
	defer apiResp.Body.Close()

	respBody, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		return nil, err
	}

	comment := &dto.Comment{}
	if err := json.Unmarshal(respBody, comment); err != nil {
		return nil, err
	}

	return &GetCommentResponse{comment}, nil
}

type UpdateCommentInput struct {
	PostID  int
	ID      int
	Content string
}

func (s *service) UpdateComment(ctx context.Context, input *UpdateCommentInput) error {
	reqInput := struct {
		Content string `json:"content"`
	}{input.Content}

	reqBody, err := json.Marshal(reqInput)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/%d/comments/%d", s.postsEndpoint, input.PostID, input.ID),
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
		return errors.New("COMMENTS: /UPDATE - API ERROR")
	}

	return nil
}

type DeleteCommentInput struct {
	PostID int
	ID     int
}

func (s *service) DeleteComment(ctx context.Context, input *DeleteCommentInput) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d/comments/%d", s.postsEndpoint, input.PostID, input.ID),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)))

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if apiResp.StatusCode != http.StatusOK {
		return errors.New("COMMENTS: /DELETE - API ERROR")
	}

	return nil
}
