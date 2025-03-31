package comments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/pkg/requests"
	"io/ioutil"
	"net/http"
)

var (
	ErrCommentsBadRequest       = errors.New("COMMENTS: bad request")
	ErrCommentsUserUnauthorized = errors.New("COMMENTS: user unauthorized")
	ErrCommentsForbidden        = errors.New("COMMENTS: forbidden")
	ErrCommentNotFound          = errors.New("DATABASE: Comment not found")
	ErrPostNotFound             = errors.New("DATABASE: Post not found")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("COMMENTS (API): %s", status)
}

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

	resp, err := requests.SendRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%d/comments", s.postsEndpoint, input.PostID),
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

	if resp.StatusCode == http.StatusBadRequest {
		// do smth
		return ErrCommentsBadRequest
	} else if resp.StatusCode == http.StatusUnauthorized {
		return ErrCommentsUserUnauthorized
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrPostNotFound
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return ErrAPIUnhandled(resp.Status)
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
	resp, err := requests.SendRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%d/comments/%d", s.postsEndpoint, input.PostID, input.ID),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrCommentNotFound
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return nil, ErrAPIUnhandled(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
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

	resp, err := requests.SendRequest(
		http.MethodPut,
		fmt.Sprintf("%s/%d/comments/%d", s.postsEndpoint, input.PostID, input.ID),
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

	if resp.StatusCode == http.StatusBadRequest {
		// do smth
		return ErrCommentsBadRequest
	} else if resp.StatusCode == http.StatusUnauthorized {
		return ErrCommentsUserUnauthorized
	} else if resp.StatusCode == http.StatusForbidden {
		return ErrCommentsForbidden
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrCommentNotFound
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return ErrAPIUnhandled(resp.Status)
	}

	return nil
}

type DeleteCommentInput struct {
	PostID int
	ID     int
}

func (s *service) DeleteComment(ctx context.Context, input *DeleteCommentInput) error {
	resp, err := requests.SendRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d/comments/%d", s.postsEndpoint, input.PostID, input.ID),
		nil,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)),
		},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrCommentsUserUnauthorized
	} else if resp.StatusCode == http.StatusForbidden {
		return ErrCommentsForbidden
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrCommentNotFound
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return ErrAPIUnhandled(resp.Status)
	}

	return nil
}
