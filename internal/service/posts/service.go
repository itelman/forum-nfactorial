package posts

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
	CreatePost(ctx context.Context, input *CreatePostInput) (*CreatePostResponse, error)
	GetPost(input *GetPostInput) (*GetPostResponse, error)
	GetAllLatestPosts() (*GetAllPostsResponse, error)
	UpdatePost(ctx context.Context, input *UpdatePostInput) error
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
	CategoriesID []string `json:"categories_id"`
}

type CreatePostResponse struct {
	PostID int `json:"id"`
}

func (s *service) CreatePost(ctx context.Context, input *CreatePostInput) (*CreatePostResponse, error) {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.postsEndpoint,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)))

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if apiResp.StatusCode != http.StatusOK {
		return nil, errors.New("POSTS: /CREATE - API ERROR")
	}
	defer apiResp.Body.Close()

	respBody, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		return nil, err
	}

	post := &dto.Post{}
	if err := json.Unmarshal(respBody, post); err != nil {
		return nil, err
	}

	return &CreatePostResponse{post.ID}, nil
}

type GetPostInput struct {
	ID int
}

type GetPostResponse struct {
	Post *dto.Post
}

func (s *service) GetPost(input *GetPostInput) (*GetPostResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%d", s.postsEndpoint, input.ID),
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
		return nil, errors.New("POSTS: /GET - API ERROR")
	}
	defer apiResp.Body.Close()

	respBody, err := ioutil.ReadAll(apiResp.Body)
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
	req, err := http.NewRequest(
		http.MethodGet,
		s.postsEndpoint,
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
		return nil, errors.New("POSTS: /LIST - API ERROR")
	}
	defer apiResp.Body.Close()

	respBody, err := ioutil.ReadAll(apiResp.Body)
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

func (s *service) UpdatePost(ctx context.Context, input *UpdatePostInput) error {
	reqInput := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{input.Title, input.Content}

	reqBody, err := json.Marshal(reqInput)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/%d", s.postsEndpoint, input.ID),
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
		return errors.New("POSTS: /UPDATE - API ERROR")
	}

	return nil
}

type DeletePostInput struct {
	ID int
}

func (s *service) DeletePost(ctx context.Context, input *DeletePostInput) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d", s.postsEndpoint, input.ID),
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
		return errors.New("POSTS: /DELETE - API ERROR")
	}

	return nil
}
