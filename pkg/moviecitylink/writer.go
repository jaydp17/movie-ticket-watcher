package moviecitylink

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

func WriteOne(dbClient dynamodbiface.ClientAPI, link MovieCityLink) error {
	links := make([]MovieCityLink, 0)
	links = append(links, link)
	return Write(dbClient, links)
}

func Write(dbClient dynamodbiface.ClientAPI, links []MovieCityLink) error {
	writables := make([]db.Writable, len(links))
	for i, l := range links {
		writables[i] = l
	}
	return db.Write(dbClient, writables, TableName)
}
