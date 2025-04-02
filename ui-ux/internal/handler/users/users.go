package users

import (
	"errors"
	"github.com/itelman/forum/internal/dto"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/service/users"
	"github.com/itelman/forum/pkg/encoder"
	"github.com/itelman/forum/pkg/templates"
	"github.com/itelman/forum/pkg/validator"
	"net/http"
)

type handlers struct {
	*handler.Handlers
	users users.Service
}

func NewHandlers(handler *handler.Handlers, users users.Service) *handlers {
	return &handlers{handler, users}
}

func (h *handlers) RegisterMux(mux *http.ServeMux) {
	authRoutes := []dto.Route{
		{"/user/signup", dto.GetPostMethods, h.signupGet},
		{"/user/login", dto.GetPostMethods, h.loginGet},
	}

	for _, route := range authRoutes {
		mux.Handle(route.Path, h.DynMiddleware.Chain(h.DynMiddleware.ForbidAuthenticatedUser(http.HandlerFunc(route.Handler)), route.Path, route.Methods))
	}

	logoutRoute := dto.Route{"/user/logout", dto.PostMethod, h.logout}
	mux.Handle(logoutRoute.Path, h.DynMiddleware.Chain(h.DynMiddleware.RequireAuthenticatedUser(http.HandlerFunc(logoutRoute.Handler)), logoutRoute.Path, logoutRoute.Methods))
}

func (h *handlers) signupGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.signupPost(w, r)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "signup_page", templates.TemplateData{
		templates.Form: validator.NewForm(nil, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) signupPost(w http.ResponseWriter, r *http.Request) {
	req, err := users.DecodeSignupUser(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*users.SignupUserInput)

	resp, err := h.users.SignupUser(input)
	if errors.Is(err, users.ErrUsersBadRequest) {
		if err := h.TmplRender.RenderData(w, r, "signup_page", templates.TemplateData{
			templates.Form: validator.NewForm(r.PostForm, resp.Errors),
		}); err != nil {
			h.Exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	h.FlashManager.UpdateFlash(dto.FlashSignupSuccessful)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *handlers) loginGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.loginPost(w, r)
		return
	}

	if err := h.TmplRender.RenderData(w, r, "login_page", templates.TemplateData{
		templates.Form: validator.NewForm(nil, nil),
	}); err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}
}

func (h *handlers) loginPost(w http.ResponseWriter, r *http.Request) {
	req, err := users.DecodeLoginUser(r)
	if err != nil {
		h.Exceptions.ErrBadRequestHandler(w, r)
		return
	}

	input := req.(*users.LoginUserInput)

	resp, err := h.users.LoginUser(input)
	if errors.Is(err, users.ErrUsersBadRequest) {
		if err := h.TmplRender.RenderData(w, r, "login_page", templates.TemplateData{
			templates.Form: validator.NewForm(r.PostForm, resp.Errors),
		}); err != nil {
			h.Exceptions.ErrInternalServerHandler(w, r, err)
			return
		}

		return
	} else if err != nil {
		h.Exceptions.ErrInternalServerHandler(w, r, err)
		return
	}

	http.SetCookie(w, dto.NewCookie(dto.TokenEncode, encoder.EncodeAccessToken(resp.AccessToken)))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handlers) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, dto.DeleteCookie(dto.TokenEncode))
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
