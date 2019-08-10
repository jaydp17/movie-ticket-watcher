package movies

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

func Write(dbClient dynamodbiface.ClientAPI, movies []Movie) error {
	writables := make([]db.Writable, len(movies))
	for i, m := range movies {
		writables[i] = m
	}
	return db.Write(dbClient, writables, TableName)
}
