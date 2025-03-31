package comments

import (
	"net/http"
	"strconv"
)

func DecodeCreateComment(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrCommentsBadRequest
	}

	postId, err := strconv.Atoi(r.PostForm.Get("post_id"))
	if err != nil {
		return nil, ErrCommentsBadRequest
	}

	return &CreateCommentInput{
		PostID:  postId,
		Content: r.PostForm.Get("content"),
	}, nil
}

func DecodeGetComment(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, ErrCommentsBadRequest
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		return nil, ErrCommentsBadRequest
	}

	return &GetCommentInput{
		PostID: postId,
		ID:     id,
	}, nil
}

func DecodeUpdateComment(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, ErrCommentsBadRequest
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		return nil, ErrCommentsBadRequest
	}

	if err := r.ParseForm(); err != nil {
		return nil, ErrCommentsBadRequest
	}

	return &UpdateCommentInput{
		PostID:  postId,
		ID:      id,
		Content: r.PostForm.Get("content"),
	}, nil
}

func DecodeDeleteComment(r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return nil, ErrCommentsBadRequest
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		return nil, ErrCommentsBadRequest
	}

	return &DeleteCommentInput{
		PostID: postId,
		ID:     id,
	}, nil
}
