// Package env provides API for handling environment variables
package env

import (
	"log"
	"os"
	"strconv"
)

// GetEnv gets value of the given environment variable name, and return default value if not exists
func GetEnv(name string, defaultValue string) string {
	value, exist := os.LookupEnv(name)
	if !exist {
		return defaultValue
	}
	return value
}

// GetEnvInt gets integer value of the given environment variable name, and return default value if not exists
func GetEnvInt(name string, defaultValue int) int {
	value, exist := os.LookupEnv(name)
	if !exist {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("invalid integer value %v of variable %v, ", value, name)
		return defaultValue
	}
	return result
}

// GetEnvInt gets integer value of the given environment variable name, and return default value if not exists
func GetEnvFloat64(name string, defaultValue float64) float64 {
	value, exist := os.LookupEnv(name)
	if !exist {
		return defaultValue
	}
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("invalid integer value %v of variable %v, ", value, name)
		return defaultValue
	}
	return result
}
