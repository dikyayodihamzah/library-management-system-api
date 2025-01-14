package userrepo

import (
	"fmt"

	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/query"
)

var SortUserMap = map[string]string{
	"name":       "full_name",
	"created_at": "created_at",
	"role":       "role",
}

func filterUsers(queryStr string, filter *model.QueryParam) (string, []interface{}) {
	if filter == nil {
		return queryStr, make([]interface{}, 0)
	}

	var args []interface{}

	if filter.Search != "" {
		queryStr = query.ClauseBuilder(queryStr) + fmt.Sprintf("full_name ILIKE $%d", len(args)+1)
		args = append(args, "%"+filter.Search+"%")
	}

	return queryStr, args
}
