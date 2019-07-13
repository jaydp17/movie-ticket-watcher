package db

import "fmt"

const dynamoTablePrefix = "ticket-watcher-"

func GetFullTableName(tableName string) string {
	return fmt.Sprintf("%s.%s", dynamoTablePrefix, tableName)
}
