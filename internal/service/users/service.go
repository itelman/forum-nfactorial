package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/itelman/forum/internal/dto"
)

type Service interface {
	SignupUser(input *SignupUserInput) error
	LoginUser(input *LoginUserInput) (*LoginUserResponse, error)
	GetAuthUser(ctx context.Context) (*GetAuthUserResponse, error)
}

type service struct {
	userEndpoint string
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
		s.userEndpoint = apiLink + "/user"
	}
}

type SignupUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *service) SignupUser(input *SignupUserInput) error {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.userEndpoint+"/signup",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("USERS: /SIGNUP - API ERROR")
	}

	return nil
}

type LoginUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	Type         string `json:"type"`
	RefreshToken string `json:"refresh_token"`
}

func (s *service) LoginUser(input *LoginUserInput) (*LoginUserResponse, error) {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.userEndpoint+"/login",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	apiResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if apiResp.StatusCode != http.StatusOK {
		return nil, errors.New("USERS: /LOGIN - API ERROR")
	}
	defer apiResp.Body.Close()

	respBody, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		return nil, err
	}

	resp := &LoginUserResponse{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type GetAuthUserInput struct {
	AccessToken string
}

type GetAuthUserResponse struct {
	User *dto.User
}

func (s *service) GetAuthUser(ctx context.Context) (*GetAuthUserResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		s.userEndpoint+"/me",
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
		return nil, errors.New("USERS: /ME - API ERROR")
	}
	defer apiResp.Body.Close()

	respBody, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		return nil, err
	}

	user := &dto.User{}
	if err := json.Unmarshal(respBody, user); err != nil {
		return nil, err
	}

	return &GetAuthUserResponse{user}, nil
}
