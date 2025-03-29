package filters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"io/ioutil"
	"net/http"
)

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
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%d/posts", s.categoriesEndpoint, input.CategoryID),
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
		return nil, errors.New("FILTERS: /{CATG_ID}/POSTS - API ERROR")
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

	return &GetPostsByCategoriesResponse{posts}, nil
}
