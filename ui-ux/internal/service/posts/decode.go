package posts

import (
	"net/http"
	"strconv"
)

func DecodeCreatePost(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrPostsBadRequest
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
		return nil, ErrPostsBadRequest
	}

	return &GetPostInput{
		ID: id,
	}, nil
}

func DecodeUpdatePost(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, ErrPostsBadRequest
	}

	if err := r.ParseForm(); err != nil {
		return nil, ErrPostsBadRequest
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
		return nil, ErrPostsBadRequest
	}

	return &DeletePostInput{
		ID: id,
	}, nil
}
