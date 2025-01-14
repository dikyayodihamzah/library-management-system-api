package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
)

// Filter generates a query string with filtering and its arguments
// It takes the base SQL query, search columns, filter parameters, and optional arguments
// This function implements subquery by default and handles various filtering scenarios
// Returns the query string and a slice of arguments
func FilterV2(sql string, searches []string, count string, param interface{}, args ...interface{}) (queryStr, queryCount string, arguments []interface{}) {
	sql = fmt.Sprintf(Processor["page"], sql)

	if query, ok := param.(Sanitize); ok {
		query.Sanitize()
	}

	params := make(map[string]interface{})
	by, err := json.Marshal(&param)
	if err != nil {
		utils.Debug(string(by), err)
	}

	if err := json.Unmarshal(by, &params); err != nil {
		utils.Debug(err)
	}

	var wheres []string
	var restricted = []string{
		"page",
		"limit",
		"search",
		"sort",
	}

	for key, val := range params {
		if strings.Contains(key, "__") {
			wheres = append(wheres, fmt.Sprint(val))
		} else if lib.NotIn(key, restricted...) && val != nil && len(fmt.Sprint(val)) > 0 {
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				wheres = append(wheres, key+" = ANY(?)")
			} else {
				wheres = append(wheres, key+" = ?")
			}
			args = append(args, val)
		}
	}

	var searching []string
	for _, search := range searches {
		if src, ok := params["search"]; ok && search != "" && len(fmt.Sprint(src)) > 0 {
			searching = append(searching, search+" ILIKE ?")
			src = "%%" + fmt.Sprint(src) + "%%"
			args = append(args, src)
		}

	}
	if len(searching) > 0 {
		wheres = append(wheres, fmt.Sprintf("(%s)", strings.Join(searching, " OR ")))
	}
	if len(wheres) > 0 {
		sql += Where(wheres...)
	}

	filter, arguments := Statement(sql, args...)

	queryCount, _ = Count(count, filter, arguments...)

	var queryParam model.QueryParam
	lib.Merge(param, &queryParam)
	utils.Debug(queryParam)
	queryStr = Page(filter, queryParam)

	defer Reset()
	return
}
