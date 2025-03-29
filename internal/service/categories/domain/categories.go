package domain

import (
	"errors"
)

var (
	ErrCategoriesBadRequest = errors.New("CATEGORIES: bad request")
	ErrCategoryNotFound     = errors.New("DATABASE: Category not found")
)
