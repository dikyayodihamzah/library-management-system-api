package query

import (
	"fmt"
	"strings"

	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
)

// Find generates a SELECT query string and its arguments
// It takes the table name, columns to select, and a condition string
// Returns the query string and a slice of arguments
func Findv2(tableName string, column []string, condition string, args ...any) (string, []any) {
	queryStr := fmt.Sprintf(Processor["find"], strings.Join(SanitizeColumns(column), ", "), tableName, condition)

	Debug(queryStr, args...)
	defer Reset()
	return Statement(queryStr, args...)
}

// CreateV2 is an improved version of Create
// It generates an INSERT query string and its arguments
// The model should implement the model.Table interface
// This function automatically calls BeforeCreate() if the model implements it
// Returns the query string and a slice of arguments
func CreateV2(models interface{}) (string, []interface{}) {
	var (
		val       []string
		columns   = PrintKey(models)
		tablename string
	)

	if v, ok := models.(BeforeCreate); ok && !SkipHook {
		v.BeforeCreate()
	}

	if t, ok := models.(model.Table); ok {
		table := t.Table()
		tablename = table.Name
	} else {
		utils.Debug("this model doesn't implement Table interface")
	}

	for i := range columns {
		val = append(val, fmt.Sprintf(`$%d`, i+1))
	}

	queryStr := fmt.Sprintf(Processor["create"], tablename,
		fmt.Sprintf(`(%s)`, strings.Join(SanitizeColumns(columns), ", ")),
		fmt.Sprintf(`(%s)`, strings.Join(val, ", ")),
	)

	args := PrintValues(models, columns)

	Debug(queryStr, args...)
	defer Reset()
	return queryStr, args
}

// CreateBatchV2 is a generic version of CreateBatch
// It generates an INSERT query string for multiple records of type T
// Note: This function doesn't implement BeforeCreate(), ensure it's called for each model if needed
// Also, it doesn't support []model as arguments, you need to loop and parse manually to []T
// Returns the query string and a slice of arguments
func CreateBatchV2[T any](tableName string, columns []string, models ...T) (string, []interface{}) {
	var all []string
	var args []interface{}
	for j := range models {

		var val []string
		for i := range columns {
			val = append(val, fmt.Sprintf(`$%d`, (j*len(columns))+(i+1)))
		}
		value := fmt.Sprintf(`(%s)`, strings.Join(val, ", "))
		all = append(all, value)
		args = append(args, PrintValues(models[j], columns)...)
	}

	queryStr := fmt.Sprintf(Processor["create"], tableName,
		fmt.Sprintf(`(%s)`, strings.Join(SanitizeColumns(columns), ", ")),
		strings.Join(all, ", "),
	)

	Debug(queryStr, args...)
	defer Reset()
	return queryStr, args
}

// Update generates an UPDATE query string and its arguments
// It takes a model to update, a condition string, and optional arguments
// This function calls BeforeUpdate() if the model implements it
// Returns the query string and a slice of arguments
// Update will set null to value if nil
func UpdateV2(mod interface{}, condition string, args ...interface{}) (string, []interface{}) {
	if v, ok := mod.(BeforeUpdate); ok {
		v.BeforeUpdate()
	}
	table := model.TableData{}
	if t, ok := mod.(model.Table); ok {
		table = t.Table()
	}
	var js = make(map[string]interface{})
	keys := PrintKey(mod)
	vals := PrintValues(mod, keys)
	keys = SanitizeColumns(keys)
	for i := range keys {
		js[keys[i]] = vals[i]
	}
	delete(js, "id")
	delete(js, "created_at")
	delete(js, "created_by")

	return Save(table.Name, js, condition, args...)
}

// Updates generates an UPDATE query string and its arguments
// It takes a model to update, a condition string, and optional arguments
// This function calls BeforeUpdate() if the model implements it
// Returns the query string and a slice of arguments
// Updates will ignore nil value to update
func UpdatesV2(mod interface{}, condition string, args ...interface{}) (string, []interface{}) {
	if v, ok := mod.(BeforeUpdate); ok {
		v.BeforeUpdate()
	}
	table := model.TableData{}
	if t, ok := mod.(model.Table); ok {
		table = t.Table()
	}
	var js = make(map[string]interface{})
	keys := PrintKey(mod)
	vals := PrintValues(mod, keys)
	keys = SanitizeColumns(keys)
	for i := range keys {
		js[keys[i]] = vals[i]
		if vals[i] == nil {
			delete(js, keys[i])
		}
	}
	delete(js, "id")
	delete(js, "created_at")
	delete(js, "created_by")

	return Save(table.Name, js, condition, args...)
}

// Save generates an UPDATE query string and its arguments
// It takes the table name, a map of data to update, a condition string, and optional arguments
// Returns the query string and a slice of arguments
func SaveV2(tableName string, data map[string]interface{}, condition string, args ...interface{}) (string, []interface{}) {
	var (
		values        []interface{}
		keys, updates []string
	)
	values = append(values, args...)
	for key, value := range data {
		if value != nil || key != "created_at" {
			keys = append(keys, key)
			values = append(values, value)
		}
	}

	for i := range keys {
		updates = append(updates, fmt.Sprintf(`%s = $%d`, keys[i], i+len(args)+1))
	}

	queryStr := fmt.Sprintf(Processor["update"], tableName, strings.Join(updates, ", "), condition)
	Debug(queryStr, values...)
	defer Reset()
	return Statement(queryStr, values...)

}

// Delete generates a soft delete UPDATE query string and its arguments
// It takes the table name, a condition string, and optional arguments
// This function sets the deleted_at column to the current timestamp
// Returns the query string and a slice of arguments
func DeleteV2(v any, condition string, args ...interface{}) (string, []interface{}) {
	now := lib.TimeNow().Format(lib.TimeFormat())
	var tableName string
	if t, ok := v.(model.Table); ok {
		tableName = t.Table().Name
	}
	queryStr := fmt.Sprintf(Processor["update"], tableName, fmt.Sprintf(`deleted_at = '%s'`, now), condition)

	Debug(queryStr, args...)
	defer Reset()
	return Statement(queryStr, args...)
}

// HardDelete generates a DELETE query string and its arguments
// It takes the table name, a condition string, and optional arguments
// This function performs a permanent deletion of records
// Returns the query string and a slice of arguments
func HardDeleteV2(v any, condition string, args ...interface{}) (string, []interface{}) {
	var tableName string
	if t, ok := v.(model.Table); ok {
		tableName = t.Table().Name
	}
	queryStr := fmt.Sprintf(Processor["delete"], tableName, condition)

	Debug(queryStr, args...)
	defer Reset()
	return Statement(queryStr, args...)
}
