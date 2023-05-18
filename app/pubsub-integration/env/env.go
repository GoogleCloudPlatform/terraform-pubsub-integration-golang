// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// GetEnvFloat64 gets integer value of the given environment variable name, and return default value if not exists
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
