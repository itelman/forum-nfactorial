package middleware

import (
	"context"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/service/users"
	"github.com/itelman/forum/pkg/encoder"
	"github.com/itelman/forum/pkg/flash"
	"net/http"
)

type AuthMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}

type middleware struct {
	users        users.Service
	exceptions   exception.Exceptions
	flashManager flash.FlashManager
}

func NewMiddleware(users users.Service, exceptions exception.Exceptions, flashManager flash.FlashManager) *middleware {
	return &middleware{
		users:        users,
		exceptions:   exceptions,
		flashManager: flashManager,
	}
}

func (m *middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errInternalSrvResp := func(err error) {
			http.SetCookie(w, dto.DeleteCookie(dto.TokenEncode))
			m.exceptions.ErrInternalServerHandler(w, r, err)
		}

		cookie, err := r.Cookie(dto.TokenEncode)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		accessToken, err := encoder.DecodeAccessToken(cookie.Value)
		if err != nil {
			errInternalSrvResp(err)
			return
		}

		ctx := context.WithValue(r.Context(), dto.ContextKeyToken, accessToken)

		resp, err := m.users.GetAuthUser(ctx)
		if errors.Is(err, users.ErrUsersUnauthorized) {
			m.flashManager.UpdateFlash(dto.FlashSessionExpired)
			http.SetCookie(w, dto.DeleteCookie(dto.TokenEncode))
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			errInternalSrvResp(err)
			return
		}

		ctx = context.WithValue(ctx, dto.ContextKeyUser, resp.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
