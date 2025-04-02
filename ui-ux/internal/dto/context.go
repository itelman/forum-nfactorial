package dto

import (
	"context"
	"net/http"
)

type contextKey string

var (
	ContextKeyUser  = contextKey("user")
	ContextKeyToken = contextKey("token")
)

func GetAuthUser(r *http.Request) *User {
	val := r.Context().Value(ContextKeyUser)

	user, ok := val.(*User)
	if !ok {
		return nil
	}

	return user
}

func GetAccessToken(ctx context.Context) string {
	val := ctx.Value(ContextKeyToken)

	token, ok := val.(string)
	if !ok {
		return ""
	}

	return token
}
