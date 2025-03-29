package comment_reactions

import (
	"github.com/itelman/forum/internal/service/comment_reactions/domain"
	"net/http"
	"strconv"
)

func DecodeCreateCommentReaction(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrCommentReactionsBadRequest
	}

	commentId, err := strconv.Atoi(r.PostForm.Get("comment_id"))
	if err != nil {
		return nil, domain.ErrCommentReactionsBadRequest
	}

	postId, err := strconv.Atoi(r.PostForm.Get("post_id"))
	if err != nil {
		return nil, domain.ErrCommentReactionsBadRequest
	}

	isLike, err := strconv.Atoi(r.PostForm.Get("is_like"))
	if err != nil {
		return nil, domain.ErrCommentReactionsBadRequest
	}

	if !(isLike == 0 || isLike == 1) {
		return nil, domain.ErrCommentReactionsBadRequest
	}

	return &CreateCommentReactionInput{
		CommentID: commentId,
		PostID:    postId,
		IsLike:    isLike,
	}, nil
}
