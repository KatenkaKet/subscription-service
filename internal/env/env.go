package env

import (
	"fmt"
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

func GetPostgresDSN() string {
	driver := GetEnvString("DB_DRIVER", "postgres")
	user := GetEnvString("DB_USER", "myuser")
	password := GetEnvString("DB_PASSWORD", "123")
	host := GetEnvString("DB_HOST", "localhost")
	port := GetEnvString("DB_PORT", "5432")
	dbName := GetEnvString("DB_NAME", "subscription_service")
	sslMode := GetEnvString("DB_SSLMODE", "disable")

	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		driver, user, password, host, port, dbName, sslMode)
}
