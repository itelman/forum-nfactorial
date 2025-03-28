package main

import (
	"github.com/itelman/forum/pkg/templates"
)

type Dependencies struct {
	templateCache templates.TemplateCache
}

func NewDependencies(opts ...Option) (deps *Dependencies, err error) {
	deps = &Dependencies{}
	for _, opt := range opts {
		if err := opt(deps); err != nil {
			return nil, err
		}
	}

	return deps, nil
}

type Option func(*Dependencies) error

func WithTemplateCache(dir string) Option {
	return func(d *Dependencies) error {
		templateCache, err := templates.NewTemplateCache(dir)
		if err != nil {
			return err
		}

		d.templateCache = templateCache
		return nil
	}
}
