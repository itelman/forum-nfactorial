package domain

import (
	"errors"
)

var (
	ErrFiltersBadRequest   = errors.New("FILTERS: bad request")
	ErrFiltersNoneSelected = errors.New("FILTERS: none selected")
	ErrUserUnauthorized    = errors.New("FILTERS: user unauthorized")
)
