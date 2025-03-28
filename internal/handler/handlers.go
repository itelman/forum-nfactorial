package handler

import (
	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/middleware/dynamic"
	"github.com/itelman/forum/pkg/templates"
)

type Handlers struct {
	DynMiddleware dynamic.DynamicMiddleware
	Exceptions    exception.Exceptions
	TmplRender    templates.TemplateRender
}

func NewHandlers(
	dynMiddleware dynamic.DynamicMiddleware,
	exceptions exception.Exceptions,
	tmplRender templates.TemplateRender,
) *Handlers {
	return &Handlers{
		DynMiddleware: dynMiddleware,
		Exceptions:    exceptions,
		TmplRender:    tmplRender,
	}
}
