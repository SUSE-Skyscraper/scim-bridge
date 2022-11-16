package pagination

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Params struct {
	Offset int32
	Limit  int32
}

var defaultPaginationParams = Params{
	Offset: 0,
	Limit:  10,
}

func Paginate(r *http.Request) Params {
	startIndex := chi.URLParam(r, "startIndex")
	count := chi.URLParam(r, "count")

	var err error
	var limit int64
	var offset int64

	if startIndex == "" {
		limit = 10
	} else {
		limit, err = strconv.ParseInt(count, 10, 32)
		if err != nil {
			return defaultPaginationParams
		}
	}

	if count == "" {
		offset = 0
	} else {
		offset, err = strconv.ParseInt(startIndex, 10, 32)
		if err != nil {
			return defaultPaginationParams
		}
	}

	params := Params{
		Offset: int32(offset),
		Limit:  int32(limit),
	}

	return params
}
