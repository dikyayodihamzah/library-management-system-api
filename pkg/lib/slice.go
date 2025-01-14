package lib

import "sort"

// FillSlice fill in slice with data
func FillSlice[T any](slice []T, value T, index int) []T {
	if index == 0 {
		var first []T
		first = append(first, value)
		first = append(first, slice...)
		return first
	}

	if index == len(slice) {
		slice = append(slice, value)
		return slice
	}

	// todo : insert in middle of array

	return slice
}

func LastLoop[T any](loop []T, i int) bool {
	return i == len(loop)-1
}

func FirstLoop(i int) bool {
	return i == 0
}

// RemoveDuplicate to remove duplicate data in slice
func RemoveDuplicate[T comparable](slice []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range slice {
		if _, ok := allKeys[item]; !ok {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// FindInSlice return true when given key found in slice
func FindInSlice[T comparable](key T, slice ...T) bool {
	for _, val := range slice {
		if key == val {
			return true
		}
	}

	return false
}

func NotIn[T comparable](key T, from ...T) bool {
	for _, s := range from {
		if key == s {
			return false
		}
	}
	return true
}

func RemoveNil[T any](v []*T) []*T {
	var new []*T
	for i := range v {
		if v[i] != nil {
			new = append(new, v[i])
		}
	}
	return new
}

func RemoveSlice[T comparable](slc []T, keys ...T) []T {
	var res []T
	for _, s := range slc {
		if !FindInSlice(s, keys...) {
			res = append(res, s)
		}
	}
	return res
}

// Pages to trim slice with given page and limit like in query db
func Pages[T any](v []T, page, limit int) []T {
	count := len(v)
	if limit > 0 {
		if page > 0 {
			page--
		}

		offset := page * limit
		if offset > count {
			return nil
		}
		if offset+limit > count {
			return v[offset:]
		}
		return v[offset : offset+limit]
	}
	return v
}

func IsValidSlicePtr[T any](v []*T) bool {
	for i := range v {
		if v[i] == nil {
			return false
		}
	}
	return true
}

type sortable interface {
	string | int | int64 | float32 | float64
}

// CompareSlice to compare 2 slice, index matter by default ([a, b] != [b, a])
// ignoreIndex to ignoring index and compare value only ([a, b] == [b, a])
func CompareSlice[T sortable](a, b []T, ignoreIndex ...bool) bool {
	if len(a) != len(b) {
		return false
	}

	if len(ignoreIndex) > 0 && ignoreIndex[0] {
		sort.Slice(a, func(i, j int) bool {
			return a[i] < a[j]
		})

		sort.Slice(b, func(i, j int) bool {
			return b[i] < b[j]
		})
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// // MakeSlice make single data to slice
// func MakeSlice[T any](v T) []T {
// 	var slice []T
// 	slice = append(slice, v)

// 	return slice
// }
