package cinemas

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

func Write(dbClient dynamodbiface.ClientAPI, cinemas []Cinema) error {
	writables := make([]db.Writable, len(cinemas))
	for i, c := range cinemas {
		writables[i] = c
	}
	return db.Write(dbClient, writables, TableName)
}
