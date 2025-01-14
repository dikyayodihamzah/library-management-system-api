package query

import "fmt"

func Paginate(queryStr string, page, limit int) string {
	if limit == 0 {
		return queryStr
	}

	if page > 0 {
		queryStr += fmt.Sprintf("\n LIMIT %d OFFSET %d", limit, (page-1)*limit)
	}

	return queryStr
}
