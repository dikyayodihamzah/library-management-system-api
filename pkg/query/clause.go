package query

import "strings"

func ClauseBuilder(query string) string {
	if strings.Contains(query, "WHERE") {
		return query + " AND "
	} else {
		return query + " WHERE "
	}
}
