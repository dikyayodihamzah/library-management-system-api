package query

import (
	"reflect"
	"time"
)

func ToQuestType(fieldType reflect.Type) string {
	switch fieldType.Kind() {
	case reflect.String:
		return "VARCHAR"
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "INT"
	case reflect.Float32, reflect.Float64:
		return "DOUBLE"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP"
		}
	}

	return "VARCHAR"
}
