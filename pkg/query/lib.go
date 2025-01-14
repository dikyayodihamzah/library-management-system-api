package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
	"github.com/google/uuid"
)

func PrintKey(any interface{}) []string {
	return printFields(any, []string{})
}

func printFields(any interface{}, k []string) (keys []string) {
	keys = append(keys, k...)
	val := reflect.ValueOf(any)

	if val.Kind() == reflect.Pointer {
		val = reflect.Indirect(val)
	}

	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name
		if t.Type.Kind() == reflect.Struct {
			keys = append(keys, printFields(val.Field(i).Interface(), k)...)
		}

		// t.Type.Elem()

		switch jsonTag := t.Tag.Get("json"); jsonTag {
		case "-", "", "deleted_at":
		default:
			parts := strings.Split(jsonTag, ",")
			name := parts[0]
			if name == "" {
				name = ToSnake(fieldName)
			}
			if dbTag := t.Tag.Get("db"); dbTag != "" {
				if dbTag == "-" {
					continue
				}
				name += "," + strings.Split(dbTag, ",")[0]
			}
			keys = append(keys, name)
		}

	}
	return
}

func PrintDBKey(any interface{}) []string {
	return printDBFields(any, []string{})
}

func printDBFields(any interface{}, k []string) (keys []string) {
	keys = append(keys, k...)
	val := reflect.ValueOf(any)

	if val.Kind() == reflect.Pointer {
		val = reflect.Indirect(val)
	}

	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name
		if t.Type.Kind() == reflect.Struct {
			keys = append(keys, printDBFields(val.Field(i).Interface(), k)...)

		}

		switch dbTag := t.Tag.Get("db"); dbTag {
		case "-", "", "deleted_at":
		default:
			parts := strings.Split(dbTag, ",")
			name := parts[0]
			if name == "" {
				name = ToSnake(fieldName)
			}
			keys = append(keys, name)
		}

	}
	return
}

func PrintValues(v interface{}, keys []string) []interface{} {
	var vals []interface{}
	by, _ := json.Marshal(v)
	var js = make(map[string]interface{})
	json.Unmarshal(by, &js)

	for _, key := range keys {
		key = strings.Split(key, ",")[0]
		val, ok := js[key]
		if ok {
			vals = append(vals, val)
		} else {
			vals = append(vals, nil)
		}
	}
	return vals
}

func Debug(str string, args ...interface{}) {
	for i := range args {
		s, _ := ToString(args[i])
		str = strings.Replace(str, "?", s, 1)
		str = strings.Replace(str, fmt.Sprintf("$%d", i+1), s, 1)
	}

	utils.Debug(str)
}

func ToString(v interface{}) (string, error) {
	valOf := reflect.ValueOf(v)
	var res string

	if u, err := uuid.Parse(fmt.Sprint(v)); err == nil {
		return "'" + u.String() + "'", nil
	}

	kind := valOf.Kind()
	if kind == reflect.Pointer {
		if valOf.IsNil() {
			return "null", nil
		} else {
			return ToString(reflect.Indirect(valOf).Interface())
		}
	}
	switch kind {
	case reflect.String:
		res = "'" + v.(string) + "'"
	case reflect.Float64, reflect.Float32, reflect.Int, reflect.Int64:
		res = fmt.Sprint(v)
	case reflect.Bool:
		res = fmt.Sprint(v)
	case reflect.Slice:
		for i := 0; i < valOf.Len(); i++ {
			var str string
			str, err := ToString(valOf.Index(i).Interface())
			if err != nil {
				return "", err
			}
			res += ", " + str
		}
		res = strings.Trim(res, ", ")
		res = "(" + res + ")"
	case reflect.Struct, reflect.Invalid:
		return "", fmt.Errorf("data type invalid")
	default:
		res = "'" + v.(string) + "'"
	}

	return res, nil
}

func ToSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

func SanitizeColumns(columns []string) []string {
	var sanitized []string
	for _, column := range columns {
		split := strings.Split(column, ",")
		newColumn := split[0]
		if len(split) > 1 {
			newColumn = split[1]
		}
		sanitized = append(sanitized, newColumn)
	}
	return sanitized
}
