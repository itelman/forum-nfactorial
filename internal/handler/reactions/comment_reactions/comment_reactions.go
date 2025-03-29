package comment_reactions

import (
	"fmt"
	"net/http"

	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/comment_reactions"
)

type commentReactionHandlers struct {
	*handler.Handlers
	commentReactions comment_reactions.Service
}

func NewHandlers(
	handler *handler.Handlers,
	commentReactions comment_reactions.Service,
) *commentReactionHandlers {
	return &commentReactionHandlers{handler, commentReactions}
}

func (h *commentReactionHandlers) RegisterMux(mux *http.ServeMux) {
	route := dto.Route{Path: "/user/posts/comments/react", Methods: dto.PostMethod, Handler: h.create}
	mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
}

func (h *commentReactionHandlers) create(w http.ResponseWriter, r *http.Request) {
	req, err := comment_reactions.DecodeCreateCommentReaction(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*comment_reactions.CreateCommentReactionInput)

	if err := h.commentReactions.CreateCommentReaction(r.Context(), input); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", input.PostID), http.StatusSeeOther)
}
