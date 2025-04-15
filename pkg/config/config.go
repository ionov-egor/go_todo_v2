package config

import (
	"os"
	"strconv"
)

// GetEnv получает переменную окружения или возвращает значение по умолчанию
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// GetEnvInt получает переменную окружения и преобразует в int
func GetEnvInt(key string, fallback int) (int, int) {
	value := GetEnv(key, "")
	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback, 0
	}
	return result, 0
}
