package lib

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// StrToInt make sting to int
func StrToInt(s string) int {
	if i, ok := strconv.Atoi(s); ok == nil {
		return i
	}
	return 0
}

// StrToUUID make string to uuid, fallback is nil
func StrToUUID(s string) *uuid.UUID {
	if u, err := uuid.Parse(s); err == nil {
		return &u
	}

	return nil
}

func StrToBool(s string) bool {
	return strings.EqualFold(s, "true")
}

// IntToRoman make roman style from given int
func IntToRoman(num int) string {
	var roman string = ""
	var numbers = []int{1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000}
	var romans = []string{"I", "IV", "V", "IX", "X", "XL", "L", "XC", "C", "CD", "D", "CM", "M"}
	var index = len(romans) - 1

	for num > 0 {
		for numbers[index] <= num {
			roman += romans[index]
			num -= numbers[index]
		}
		index -= 1
	}

	return roman
}

// Swap value from a to b and b to a
func Swap[T any](a, b *T) {
	var x T = *a
	*a = *b
	*b = x
}

// Operator will return ok value if condition is true, and not ok if condition false
//
// this func is used to fill data within one line
func Operator[T any](cond bool, ok, notOk T) T {
	if cond {
		return ok
	} else {
		return notOk
	}
}


