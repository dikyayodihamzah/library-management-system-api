package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
)

// Find generates a SELECT query string and its arguments
// It takes the table name, columns to select, and a condition string
// Returns the query string and a slice of arguments
func Find(tableName string, column []string, condition string, args ...any) (string, []any) {
	queryStr := fmt.Sprintf(Processor["find"], strings.Join(SanitizeColumns(column), ", "), tableName, condition)

	Debug(queryStr, args...)
	defer Reset()
	return Statement(queryStr, args...)
}

// Create generates an INSERT query string and its arguments
// It takes the table name, column names, and a model to insert
// This function automatically calls BeforeCreate() if the model implements it
// Returns the query string and a slice of arguments
func Create(tablename string, columns []string, models interface{}) (string, []interface{}) {
	var val []string

	if v, ok := models.(BeforeCreate); ok && !SkipHook {
		v.BeforeCreate()
	}

	table := model.TableData{}
	if t, ok := models.(model.Table); ok {
		table = t.Table()
	}

	for i := range columns {
		val = append(val, fmt.Sprintf(`$%d`, i+1))
	}

	if tablename == "" {
		tablename = table.Name
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

// CreateBatch generates an INSERT query string for multiple records
// It takes the table name, column names, and a variadic number of models to insert
// Note: This function doesn't implement BeforeCreate(), ensure it's called for each model if needed
// Also, it doesn't support []model as arguments, you need to loop and parse manually to []interface{}
// Returns the query string and a slice of arguments
func CreateBatch(tableName string, columns []string, models ...interface{}) (string, []interface{}) {
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
	utils.Debug(queryStr)
	Debug(queryStr, args...)
	defer Reset()
	return queryStr, args
}

// Update generates an UPDATE query string and its arguments
// It takes a model to update, a condition string, and optional arguments
// This function calls BeforeUpdate() if the model implements it
// Returns the query string and a slice of arguments
// Update will set null to value if nil
func Update(mod interface{}, condition string, args ...interface{}) (string, []interface{}) {
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
func Updates(mod interface{}, condition string, args ...interface{}) (string, []interface{}) {
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
func Save(tableName string, data map[string]interface{}, condition string, args ...interface{}) (string, []interface{}) {
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
func Delete(tablename string, condition string, args ...interface{}) (string, []interface{}) {
	now := lib.TimeNow().Format(lib.TimeFormat())
	queryStr := fmt.Sprintf(Processor["update"], tablename, fmt.Sprintf(`deleted_at = '%s'`, now), condition)

	Debug(queryStr, args...)
	defer Reset()
	return Statement(queryStr, args...)
}

// HardDelete generates a DELETE query string and its arguments
// It takes the table name, a condition string, and optional arguments
// This function performs a permanent deletion of records
// Returns the query string and a slice of arguments
func HardDelete(tableName string, condition string, args ...interface{}) (string, []interface{}) {
	queryStr := fmt.Sprintf(Processor["delete"], tableName, condition)

	Debug(queryStr, args...)
	defer Reset()
	return Statement(queryStr, args...)
}

// Count generates a COUNT query string and its arguments
// It takes the column to count, the base SQL query, and optional arguments
// Returns the query string and a slice of arguments
func Count(what, sql string, args ...interface{}) (string, []interface{}) {
	if what == "" {
		what = "*"
	}
	var queryStr string
	if strings.Contains(sql, "SELECT * ") {
		queryStr = strings.Replace(sql, "*", fmt.Sprintf("COUNT(%s)", what), 1)
	} else {
		queryStr = fmt.Sprintf(Processor["count"], what, sql)
	}
	Debug(queryStr, args...)
	defer Reset()
	return Statement(queryStr, args...)
}

// Filter generates a query string with filtering and its arguments
// It takes the base SQL query, search columns, filter parameters, and optional arguments
// This function implements subquery by default and handles various filtering scenarios
// Returns the query string and a slice of arguments
func Filter(sql string, searches []string, param interface{}, args ...interface{}) (string, []interface{}) {
	utils.Debug("DisableSubQueries", DisableSubQueries)
	if !DisableSubQueries {
		sql = fmt.Sprintf(Processor["page"], sql)
	}

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

	defer Reset()
	return Statement(sql, args...)
}

// Page adds sorting, offset, and limit clauses to the SQL query
// It takes the base SQL query and query parameters
// Returns the modified SQL query string
func Page(sql string, param model.QueryParam) string {
	if !lib.IsEmptyString(param.Sort) {
		str := strings.Trim(param.Sort, "-")
		sort := "ASC"
		if str != param.Sort {
			sort = "DESC"
		}
		var sortValue []string
		for _, value := range strings.Split(str, ",") {
			sortValue = append(sortValue, fmt.Sprintf("%s %s", value, sort))
		}
		sql += " ORDER BY " + strings.Join(sortValue, ", ")

	}

	if param.Limit != 0 {
		if param.Page > 0 {
			param.Page--
		}
		sql += fmt.Sprintf("\n LIMIT %d OFFSET %d ", param.Limit, param.Page*param.Limit)
	}

	return sql
}

func PageQuest(sql string, param model.QueryParam) string {
	if !lib.IsEmptyString(param.Sort) {
		str := strings.Trim(param.Sort, "-")
		sort := "ASC"
		if str != param.Sort {
			sort = "DESC"
		}
		var sortValue []string
		for _, value := range strings.Split(str, ",") {
			sortValue = append(sortValue, fmt.Sprintf(`"%s" %s`, value, sort))
		}
		sql += " ORDER BY " + strings.Join(sortValue, ", ")

	}

	if param.Limit != 0 {
		sql += fmt.Sprintf("\n LIMIT %d,%d ", (param.Page-1)*param.Limit, param.Page*param.Limit)
	}

	return sql
}

// CreateQuestTable generates a CREATE TABLE query string for Quest
// It takes the table name and a model struct
// The function uses reflection to analyze the model and create appropriate columns
// Returns the CREATE TABLE query string
func CreateQuestTable(tableName string, model interface{}, designatedTimestamp string) string {
	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Struct {
		utils.Debug("model must be a struct")
		return ""
	}

	var columns []string

	// Helper function to recursively process struct fields
	var processStruct func(reflect.Type)
	processStruct = func(typ reflect.Type) {
		// Resolve pointers to their underlying types
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		// Only process struct types
		if typ.Kind() != reflect.Struct {
			return
		}

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}

			// Handle embedded structs recursively
			if field.Anonymous || field.Type.Kind() == reflect.Struct {
				// Special case: Skip handling time.Time as a struct
				if field.Type == reflect.TypeOf(time.Time{}) {
					// Add the time.Time field as a TIMESTAMP column
					columnName := field.Tag.Get("db")
					if columnName != "" {
						columnName = fmt.Sprintf(`"%s"`, columnName)
						columns = append(columns, fmt.Sprintf("%s TIMESTAMP", columnName))
					}
					continue
				}
				processStruct(field.Type)
				continue
			}

			// Skip fields without a "db" tag
			columnName := field.Tag.Get("db")
			if columnName == "" {
				continue
			}

			// Get the corresponding SQL type
			columnType := ToQuestType(field.Type)

			// Encapsulate column name with double quotes
			columnName = fmt.Sprintf(`"%s"`, columnName)

			columns = append(columns, fmt.Sprintf("%s %s", columnName, columnType))
		}
	}

	// Process the top-level struct
	processStruct(modelType)

	// Construct the CREATE TABLE query
	queryStr := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (%s) TIMESTAMP(%s)`,
		tableName, strings.Join(columns, ", "), designatedTimestamp,
	)

	return queryStr
}
