package filters

import (
	"net/http"
	"strconv"
)

func DecodeGetPostsByCategories(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrFiltersBadRequest
	}

	catgId, err := strconv.Atoi(r.PostForm.Get("category_id"))
	if err != nil {
		return nil, ErrFiltersBadRequest
	}

	return &GetPostsByCategoriesInput{
		CategoryID: catgId,
	}, nil
}
