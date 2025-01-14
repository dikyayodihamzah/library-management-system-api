package lib

import (
	"fmt"
	"strings"
)

// IsEmptyString check if one of given string is empty
func IsEmptyString(str ...string) bool {
	for _, s := range str {
		if s == "" || len([]rune(s)) == 0 {
			return true
		}
	}
	return false
}

// IsValidString check if all of given string is valid
func IsValidString(str ...string) bool {
	for _, s := range str {
		if s == "" || len([]rune(s)) == 0 {
			return false
		}
	}
	return true
}

// AppendStr append str with separator
func AppendStr(sep string, str ...string) string {
	return strings.Join(str, sep)
}

// StrAbbr make abbrebration of complex string
func StrAbbr(s string) string {
	sl := strings.Split(s, " ")

	r := ""
	for _, slc := range sl {
		f := slc[0]
		r += string(f)
	}

	return strings.ToUpper(r)
}

// TrimSpace remove whitespace and double space
func TrimSpace(str string) string {
	split := strings.Split(str, " ")

	var res []string
	for _, s := range split {
		if trimmed := strings.TrimSpace(s); IsValidString(trimmed) {
			res = append(res, trimmed)
		}
	}
	return strings.Join(res, " ")
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

func CreateNameCompany(companyName, companyLine string, branchName ...string) string {
	name := companyLine + ", " + companyLine
	if len(branchName) > 0 {
		name += fmt.Sprintf(" (%s)", branchName[0])
	}

	return name
}
