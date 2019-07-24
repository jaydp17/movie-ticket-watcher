package config

import (
	"fmt"
	"os"
)

var AppEnv string
var IsProd bool
var dynamoTablePrefix string

func init() {
	AppEnv = getEnvWithDefault("APP_ENV", "development")
	IsProd = AppEnv == "production"
	dynamoTablePrefix = fmt.Sprintf("ticket-watcher-%s", AppEnv)
}

func getEnvWithDefault(envName string, fallback string) string {
	actualValue := os.Getenv(envName)
	if len(actualValue) == 0 {
		return fallback
	}
	return actualValue
}

func FullTableName(tableName string) string {
	const tablePrefix = "ticket-watcher"
	return fmt.Sprintf("%s-%s.%s", tablePrefix, AppEnv, tableName)
}
