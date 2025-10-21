// Package env for managing env variables
package env

import (
	"os"
	"strconv"
)

func GetString(envKey string, defaultVal string) string {
	value, exists := os.LookupEnv(envKey)

	if !exists || (value == "") {
		return defaultVal
	}

	return value
}

func GetInt(envKey string, defaultVal int) int {
	value, exists := os.LookupEnv(envKey)

	if !exists || value == "" {
		return defaultVal
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal
	}

	return intValue
}

func GetBool(envKey string, defaultVal bool) bool {
	value, exists := os.LookupEnv(envKey)

	if !exists || value == "" {
		return defaultVal
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}

	return boolValue
}
