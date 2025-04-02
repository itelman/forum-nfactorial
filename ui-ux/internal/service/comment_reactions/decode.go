package comment_reactions

import (
	"net/http"
	"strconv"
)

func DecodeCreateCommentReaction(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrCommentReactionsBadRequest
	}

	commentId, err := strconv.Atoi(r.PostForm.Get("comment_id"))
	if err != nil {
		return nil, ErrCommentReactionsBadRequest
	}

	postId, err := strconv.Atoi(r.PostForm.Get("post_id"))
	if err != nil {
		return nil, ErrCommentReactionsBadRequest
	}

	isLike, err := strconv.Atoi(r.PostForm.Get("is_like"))
	if err != nil {
		return nil, ErrCommentReactionsBadRequest
	}

	if !(isLike == 0 || isLike == 1) {
		return nil, ErrCommentReactionsBadRequest
	}

	return &CreateCommentReactionInput{
		CommentID: commentId,
		PostID:    postId,
		IsLike:    isLike,
	}, nil
}
