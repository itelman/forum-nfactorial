package home

import (
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/categories"
	"github.com/itelman/forum/internal/service/filters"
	"github.com/itelman/forum/internal/service/posts"
	"github.com/itelman/forum/pkg/templates"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
)

type handlers struct {
	*handler.Handlers
	posts      posts.Service
	categories categories.Service
	filters    filters.Service
}

func NewHandlers(handler *handler.Handlers, posts posts.Service, categories categories.Service, filters filters.Service) *handlers {
	return &handlers{handler, posts, categories, filters}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	routes := []dto.Route{
		{"/", dto.GetMethod, h.home},
		{"/results", dto.PostMethod, h.results},
	}

	for _, route := range routes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(http.HandlerFunc(route.Handler), route.Path, route.Methods))
	}
}

func (h *handlers) home(w http.ResponseWriter, r *http.Request) {
	postsResp, err := h.posts.GetAllLatestPosts()
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	catgRsp, err := h.categories.GetAllCategories()
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "home_page", templates.TemplateData{
		templates.Posts:      postsResp.Posts,
		templates.Categories: catgRsp.Categories,
		templates.Form:       validator.NewForm(nil, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) results(w http.ResponseWriter, r *http.Request) {
	req, err := filters.DecodeGetPostsByCategories(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*filters.GetPostsByCategoriesInput)

	filtersResp, err := h.filters.GetPostsByCategories(input)
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	catgRsp, err := h.categories.GetAllCategories()
	if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "home_page", templates.TemplateData{
		templates.Posts:      filtersResp.Posts,
		templates.Categories: catgRsp.Categories,
		templates.Form:       validator.NewForm(nil, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}
