package subscriptions

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

func WriteOne(dbClient dynamodbiface.ClientAPI, s Subscription) error {
	subscriptions := make([]Subscription, 1)
	subscriptions[0] = s
	return Write(dbClient, subscriptions)
}

func Write(dbClient dynamodbiface.ClientAPI, subscriptions []Subscription) error {
	writables := make([]db.Writable, len(subscriptions))
	for i, s := range subscriptions {
		writables[i] = s
	}
	return db.Write(dbClient, writables, TableName)
}
