package activity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/itelman/forum/internal/dto"
)

type Service interface {
	GetAllCreatedPosts(ctx context.Context) (*GetAllCreatedPostsResponse, error)
	GetAllReactedPosts(ctx context.Context) (*GetAllReactedPostsResponse, error)
}

type service struct {
	userPostsEndpoint string
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
		s.userPostsEndpoint = apiLink + "/user/posts"
	}
}

type GetAllCreatedPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllCreatedPosts(ctx context.Context) (*GetAllCreatedPostsResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		s.userPostsEndpoint+"/created",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)))

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if apiResp.StatusCode != http.StatusOK {
		return nil, errors.New("ACTIVITY: /USER/POSTS/CREATED - API ERROR")
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

	return &GetAllCreatedPostsResponse{posts}, nil
}

type GetAllReactedPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllReactedPosts(ctx context.Context) (*GetAllReactedPostsResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		s.userPostsEndpoint+"/reacted",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)))

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if apiResp.StatusCode != http.StatusOK {
		return nil, errors.New("ACTIVITY: /USER/POSTS/REACTED - API ERROR")
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

	return &GetAllReactedPostsResponse{posts}, nil
}
