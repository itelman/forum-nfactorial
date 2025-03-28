package categories

import (
	"encoding/json"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"io/ioutil"
	"net/http"
)

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
	req, err := http.NewRequest(
		http.MethodGet,
		s.categoriesEndpoint,
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
		return nil, errors.New("CATEGORIES: /LIST - API ERROR")
	}
	defer apiResp.Body.Close()

	respBody, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		return nil, err
	}

	var categories []*dto.Category
	if err := json.Unmarshal(respBody, &categories); err != nil {
		return nil, err
	}

	return &GetAllCategoriesResponse{categories}, nil
}
