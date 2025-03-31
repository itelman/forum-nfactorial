package filters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/pkg/requests"
	"io/ioutil"
	"net/http"
)

var (
	ErrFiltersBadRequest = errors.New("FILTERS: bad request")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("FILTERS (API): %s", status)
}

type Service interface {
	GetPostsByCategories(input *GetPostsByCategoriesInput) (*GetPostsByCategoriesResponse, error)
}

type service struct {
	categoriesEndpoint string
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
		s.categoriesEndpoint = apiLink + "/categories"
	}
}

type GetPostsByCategoriesInput struct {
	CategoryID int
}

type GetPostsByCategoriesResponse struct {
	Posts []*dto.Post
}

func (s *service) GetPostsByCategories(input *GetPostsByCategoriesInput) (*GetPostsByCategoriesResponse, error) {
	resp, err := requests.SendRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%d/posts", s.categoriesEndpoint, input.CategoryID),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrFiltersBadRequest
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

	return &GetPostsByCategoriesResponse{posts}, nil
}
