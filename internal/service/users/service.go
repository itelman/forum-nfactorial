package users

import (
	"bytes"
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
	ErrUsersBadRequest    = errors.New("USERS: bad request")
	ErrUsersUnauthorized  = errors.New("USERS: unauthorized")
	ErrInvalidCredentials = errors.New("DATABASE: Invalid credentials")
)

func ErrAPIUnhandled(status string) error {
	return fmt.Errorf("USERS (API): %s", status)
}

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

	resp, err := requests.SendRequest(
		http.MethodPost,
		s.userEndpoint+"/signup",
		bytes.NewBuffer(reqBody),
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		// do smth
		return ErrUsersBadRequest
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return ErrAPIUnhandled(resp.Status)
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

	apiResp, err := requests.SendRequest(
		http.MethodPost,
		s.userEndpoint+"/login",
		bytes.NewBuffer(reqBody),
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return nil, err
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode == http.StatusBadRequest {
		// do smth
		return nil, ErrUsersBadRequest
	} else if apiResp.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidCredentials
	} else if !(apiResp.StatusCode >= http.StatusOK && apiResp.StatusCode < http.StatusMultipleChoices) {
		return nil, ErrAPIUnhandled(apiResp.Status)
	}

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

type GetAuthUserResponse struct {
	User *dto.User
}

func (s *service) GetAuthUser(ctx context.Context) (*GetAuthUserResponse, error) {
	resp, err := requests.SendRequest(
		http.MethodGet,
		s.userEndpoint+"/me",
		nil,
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", dto.GetAccessToken(ctx))},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUsersUnauthorized
	} else if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) {
		return nil, ErrAPIUnhandled(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	user := &dto.User{}
	if err := json.Unmarshal(respBody, user); err != nil {
		return nil, err
	}

	return &GetAuthUserResponse{user}, nil
}
