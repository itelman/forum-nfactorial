package filters

import (
	"github.com/itelman/forum/internal/service/filters/domain"
	"net/http"
	"strconv"
)

func DecodeGetPostsByCategories(r *http.Request) (interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, domain.ErrFiltersBadRequest
	}

	catgId, err := strconv.Atoi(r.PostForm.Get("category_id"))
	if err != nil {
		catgId = -1
	}

	return &GetPostsByCategoriesInput{
		CategoryID: catgId,
	}, nil
}
