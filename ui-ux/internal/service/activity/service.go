package activity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/pkg/requests"
	"io/ioutil"
	"net/http"

	"github.com/itelman/forum/internal/dto"
)

var (
	ErrActivityUserUnauthorized = errors.New("ACTIVITY: user unauthorized")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("ACTIVITY (API): %s", status)
}

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
	resp, err := requests.SendRequest(
		http.MethodGet,
		s.userPostsEndpoint+"/created",
		nil,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)),
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrActivityUserUnauthorized
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
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

	return &GetAllCreatedPostsResponse{posts}, nil
}

type GetAllReactedPostsResponse struct {
	Posts []*dto.Post
}

func (s *service) GetAllReactedPosts(ctx context.Context) (*GetAllReactedPostsResponse, error) {
	resp, err := requests.SendRequest(
		http.MethodGet,
		s.userPostsEndpoint+"/reacted",
		nil,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx)),
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrActivityUserUnauthorized
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
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

	return &GetAllReactedPostsResponse{posts}, nil
}
