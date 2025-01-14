package utils

import (
	"github.com/spf13/cast"
)

func init() {
	LoadEnv()
}

func GetString(key string, defaultValue ...string) string {
	if val, ok := envMap[key]; ok {
		return cast.ToString(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

func GetInt(key string, defaultValue ...int) int {
	if val, ok := envMap[key]; ok {
		return cast.ToInt(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

func GetBool(key string, defaultValue ...bool) bool {
	if val, ok := envMap[key]; ok {
		return cast.ToBool(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return false
}
