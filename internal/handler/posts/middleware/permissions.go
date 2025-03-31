package middleware

import (
	"context"
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/service/posts"
	"net/http"
)

type contextKey string

var (
	ContextKeyPost = contextKey("post")
)

type PostsCheckPermissionMiddleware interface {
	CheckUserPermissions(next http.Handler) http.Handler
}

type middleware struct {
	posts      posts.Service
	exceptions exception.Exceptions
}

func NewMiddleware(posts posts.Service, exceptions exception.Exceptions) *middleware {
	return &middleware{
		posts:      posts,
		exceptions: exceptions,
	}
}

func (m *middleware) CheckUserPermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postReq, err := posts.DecodeGetPost(r)
		if err != nil {
			m.exceptions.ErrBadRequestHandler(w, r)
			return
		}

		input := postReq.(*posts.GetPostInput)

		postResp, err := m.posts.GetPost(r.Context(), input)
		if errors.Is(err, posts.ErrPostNotFound) {
			m.exceptions.ErrNotFoundHandler(w, r)
			return
		} else if err != nil {
			m.exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		if postResp.Post.User.ID != dto.GetAuthUser(r).ID {
			m.exceptions.ErrForbiddenHandler(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyPost, postResp.Post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPostFromContext(r *http.Request) *dto.Post {
	val := r.Context().Value(ContextKeyPost)

	post, ok := val.(*dto.Post)
	if !ok {
		return nil
	}

	return post
}
