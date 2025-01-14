package query

import (
	"fmt"
	"strings"
)

func Where(str ...string) string {
	wh := "WHERE "
	if DisableSubQueries {
		wh = " AND "
	}
	return wh + strings.Join(str, " AND ")
}

func NotDeletedAndStatusActive(prefix ...string) string {
	if len(prefix) == 0 {
		return "deleted_at is null AND status = 1"
	}

	var where []string
	for i := range prefix {
		var pref string
		if len(prefix[i]) > 0 {
			pref = prefix[i] + "."
		}
		where = append(where, fmt.Sprintf("%sdeleted_at is null AND %sstatus = 1", pref, pref))
	}

	return strings.Join(where, " AND ")
}

func NotDeleted(prefix ...string) string {
	if len(prefix) == 0 {
		return "deleted_at is null"
	}

	var where []string
	for i := range prefix {
		var pref string
		if len(prefix[i]) > 0 {
			pref = prefix[i] + "."
		}
		where = append(where, fmt.Sprintf("%sdeleted_at is null", pref))
	}

	return strings.Join(where, " AND ")
}
