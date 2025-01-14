package query

import (
	"fmt"
	"strings"
)

var (
	SkipHook          bool = false
	Error             error
	DisableSubQueries bool = false
)

type BeforeCreate interface {
	BeforeCreate()
}

type BeforeUpdate interface {
	BeforeUpdate()
}

type Sanitize interface {
	Sanitize()
}

var Processor map[string]string = map[string]string{
	"find":   "SELECT %s FROM %s %s",
	"create": "INSERT INTO %s %s VALUES %s",
	"update": "UPDATE %s SET %s %s",
	"delete": "DELETE FROM %s %s",
	"page":   "SELECT * FROM (%s) AS s ",
	"count":  "SELECT COUNT(%s) FROM (%s) AS c ",
}

func DisableSubQuery() {
	DisableSubQueries = true
}

func SkipHooks() {
	SkipHook = true
}

func Reset() {
	SkipHook = false
	DisableSubQueries = false
}

func Build(str string, args ...interface{}) string {
	return fmt.Sprintf(str, args...)
}

func Statement(str string, args ...interface{}) (string, []interface{}) {
	if len(args) == 0 {
		return str, args
	}

	for i := range args {
		str = strings.Replace(str, "?", "$"+fmt.Sprint(i+1), 1)
	}
	return str, args
}
