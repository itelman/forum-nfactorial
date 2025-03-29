package posts

import (
	"github.com/itelman/forum/internal/service/posts/domain"
	"net/http"
	"strconv"
)

func DecodeCreatePost(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	return &CreatePostInput{
		Title:        r.PostForm.Get("title"),
		Content:      r.PostForm.Get("content"),
		CategoriesID: r.PostForm["categories_id"],
	}, nil
}

func DecodeGetPost(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	return &GetPostInput{
		ID: id,
	}, nil
}

func DecodeUpdatePost(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	return &UpdatePostInput{
		ID:      id,
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}, nil
}

func DecodeDeletePost(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, domain.ErrPostsBadRequest
	}

	return &DeletePostInput{
		ID: id,
	}, nil
}
