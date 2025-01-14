package query

import (
	"fmt"
	"strings"
)

func ValidateSort(sort string, sortMap map[string]string) (key, value string, err error) {
	// check sort direction by first character
	if strings.HasPrefix(sort, "-") {
		key = strings.TrimPrefix(sort, "-")
		value = "desc"
	} else if strings.HasPrefix(sort, "+") {
		key = strings.TrimPrefix(sort, "+")
		value = "asc"
	} else {
		key = sort
		value = "asc"
	}

	// validate sort key
	if _, ok := sortMap[key]; !ok {
		var response string
		var length int
		for key := range sortMap {
			response += fmt.Sprintf("'%s'", key)
			if length != len(sortMap)-1 {
				response += ", "
			}
			length++
		}

		return "", "", fmt.Errorf("Invalid sort key. Available sort keys: %s", response)
	}

	// validate sort value
	if value != "asc" && value != "desc" {
		return "", "", fmt.Errorf("Invalid sort value. Available sort values: 'asc', 'desc'")
	}

	return key, value, nil
}

func Sort(queryStr, sort string, sortMap map[string]string) (string, error) {
	// validate sort
	if sort == "" {
		return queryStr, nil
	}

	key, value, err := ValidateSort(sort, sortMap)
	if err != nil {
		return "", err
	}

	// switch sort value if sort key is "no"
	if key == "no" {
		if value == "asc" {
			value = "desc"
		} else {
			value = "asc"
		}
	}

	queryStr += fmt.Sprintf(" ORDER BY %s %s", sortMap[key], value)

	return queryStr, nil
}
