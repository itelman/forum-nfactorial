package posts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/pkg/requests"
	"github.com/itelman/forum/pkg/validator"
	"io/ioutil"
	"net/http"
)

var (
	ErrPostsBadRequest       = errors.New("POSTS: bad request")
	ErrPostsUserUnauthorized = errors.New("POSTS: user unauthorized")
	ErrPostsForbidden        = errors.New("POSTS: forbidden")
	ErrPostNotFound          = errors.New("DATABASE: Post not found")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("POSTS (API): %s", status)
}

type Service interface {
	CreatePost(ctx context.Context, input *CreatePostInput) (*CreatePostResponse, error)
	GetPost(ctx context.Context, input *GetPostInput) (*GetPostResponse, error)
	GetAllLatestPosts() (*GetAllPostsResponse, error)
	UpdatePost(ctx context.Context, input *UpdatePostInput) (*UpdatePostResponse, error)
	DeletePost(ctx context.Context, input *DeletePostInput) error
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

type CreatePostInput struct {
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	CategoriesID []string `json:"categories"`
}

type CreatePostResponse struct {
	PostID int
	Errors validator.Errors
}

func (s *service) CreatePost(ctx context.Context, input *CreatePostInput) (*CreatePostResponse, error) {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resp, err := requests.SendRequest(
		http.MethodPost,
		s.postsEndpoint,
		bytes.NewBuffer(reqBody),
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)),
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var errs validator.Errors
		if err := json.Unmarshal(respBody, &errs); err != nil {
			return nil, err
		}

		return &CreatePostResponse{-1, errs}, ErrPostsBadRequest
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrPostsUserUnauthorized
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return nil, ErrAPIUnhandled(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	post := &dto.Post{}
	if err := json.Unmarshal(respBody, post); err != nil {
		return nil, err
	}

	return &CreatePostResponse{post.ID, nil}, nil
}

type GetPostInput struct {
	ID int
}

type GetPostResponse struct {
	Post *dto.Post
}

func (s *service) GetPost(ctx context.Context, input *GetPostInput) (*GetPostResponse, error) {
	var headers map[string]string = nil
	if dto.GetAccessToken(ctx) != "" {
		headers = map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)),
		}
	}

	resp, err := requests.SendRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%d", s.postsEndpoint, input.ID),
		nil,
		headers,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrPostNotFound
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return nil, ErrAPIUnhandled(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	post := &dto.Post{}
	if err := json.Unmarshal(respBody, post); err != nil {
		return nil, err
	}

	return &GetPostResponse{post}, nil
}

type GetAllPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllLatestPosts() (*GetAllPostsResponse, error) {
	resp, err := requests.SendRequest(
		http.MethodGet,
		s.postsEndpoint,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return nil, ErrAPIUnhandled(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var posts []*dto.Post
	if err := json.Unmarshal(respBody, &posts); err != nil {
		return nil, err
	}

	return &GetAllPostsResponse{posts}, nil
}

type UpdatePostInput struct {
	ID      int
	Title   string
	Content string
}

type UpdatePostResponse struct {
	Errors validator.Errors
}

func (s *service) UpdatePost(ctx context.Context, input *UpdatePostInput) (*UpdatePostResponse, error) {
	reqInput := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{input.Title, input.Content}

	reqBody, err := json.Marshal(reqInput)
	if err != nil {
		return nil, err
	}

	resp, err := requests.SendRequest(
		http.MethodPut,
		fmt.Sprintf("%s/%d", s.postsEndpoint, input.ID),
		bytes.NewBuffer(reqBody),
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)),
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var errs validator.Errors
		if err := json.Unmarshal(respBody, &errs); err != nil {
			return nil, err
		}

		return &UpdatePostResponse{errs}, ErrPostsBadRequest
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrPostsUserUnauthorized
	} else if resp.StatusCode == http.StatusForbidden {
		return nil, ErrPostsForbidden
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, ErrPostNotFound
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return nil, ErrAPIUnhandled(resp.Status)
	}

	return nil, nil
}

type DeletePostInput struct {
	ID int
}

func (s *service) DeletePost(ctx context.Context, input *DeletePostInput) error {
	resp, err := requests.SendRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d", s.postsEndpoint, input.ID),
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
		return ErrPostsUserUnauthorized
	} else if resp.StatusCode == http.StatusForbidden {
		return ErrPostsForbidden
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrPostNotFound
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return ErrAPIUnhandled(resp.Status)
	}

	return nil
}
