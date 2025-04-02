package post_reactions

import (
	"net/http"
	"strconv"
)

func DecodeCreatePostReaction(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrPostReactionsBadRequest
	}

	postId, err := strconv.Atoi(r.PostForm.Get("post_id"))
	if err != nil {
		return nil, ErrPostReactionsBadRequest
	}

	isLike, err := strconv.Atoi(r.PostForm.Get("is_like"))
	if err != nil {
		return nil, ErrPostReactionsBadRequest
	}

	if !(isLike == 0 || isLike == 1) {
		return nil, ErrPostReactionsBadRequest
	}

	return &CreatePostReactionInput{
		PostID: postId,
		IsLike: isLike,
	}, nil
}
