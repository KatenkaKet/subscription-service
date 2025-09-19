package env

import (
	"os"
	"strconv"
)

func GetEnvString(key, defaultvalue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultvalue
}

func GetEnvInt(key string, defaultvalue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultvalue
}
