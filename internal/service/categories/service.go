package categories

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
	ErrCategoriesBadRequest = errors.New("CATEGORIES: bad request")
	ErrCategoryNotFound     = errors.New("DATABASE: Category not found")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("CATEGORIES (API): %s", status)
}

type Service interface {
	GetAllCategories() (*GetAllCategoriesResponse, error)
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

type GetAllCategoriesResponse struct {
	Categories []*dto.Category
}

func (s *service) GetAllCategories() (*GetAllCategoriesResponse, error) {
	resp, err := requests.SendRequest(
		http.MethodGet,
		s.categoriesEndpoint,
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

	var categories []*dto.Category
	if err := json.Unmarshal(respBody, &categories); err != nil {
		return nil, err
	}

	return &GetAllCategoriesResponse{categories}, nil
}
