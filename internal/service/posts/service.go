package posts

import (
	"encoding/json"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"io/ioutil"
	"net/http"
)

type Service interface {
	CreatePost(input *CreatePostInput, dir string) (*CreatePostResponse, error)
	GetPost(input *GetPostInput) (*GetPostResponse, error)
	GetAllLatestPosts() (*GetAllPostsResponse, error)
	UpdatePost(input *UpdatePostInput, post *dto.Post) error
	DeletePost(input *DeletePostInput, dir string) error
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

type CreatePostResponse struct {
	PostID int
}

func (s *service) CreatePost(input *CreatePostInput, dir string) (*CreatePostResponse, error) {
	return nil, nil
}

type GetPostResponse struct {
	Post *dto.Post
}

func (s *service) GetPost(input *GetPostInput) (*GetPostResponse, error) {
	return nil, nil
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

func (s *service) UpdatePost(input *UpdatePostInput, post *dto.Post) error {
	return nil
}

func (s *service) DeletePost(input *DeletePostInput, dir string) error {
	return nil
}
