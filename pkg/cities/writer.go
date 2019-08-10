package cities

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

func Write(dbClient dynamodbiface.ClientAPI, cities []City) error {
	writables := make([]db.Writable, len(cities))
	for i, city := range cities {
		writables[i] = city
	}
	return db.Write(dbClient, writables, TableName)
}
