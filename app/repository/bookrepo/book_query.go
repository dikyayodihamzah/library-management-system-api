package bookrepo

import (
	"fmt"

	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
)

var SortBookMap = map[string]string{
	"title":      "title",
	"author":     "author",
	"genre":      "genre",
	"rating":     "rating",
	"price":      "price_idr",
	"created_at": "created_at",
}

func filterBooks(queryStr string, filter *model.QueryParam) (string, []interface{}) {
	if filter == nil {
		return queryStr, make([]interface{}, 0)
	}

	var args []interface{}

	if filter.Search != "" {
		queryStr = query.ClauseBuilder(queryStr) + fmt.Sprintf("title ILIKE $%d", len(args)+1)
		args = append(args, "%"+filter.Search+"%")
	}

	return queryStr, args
}
