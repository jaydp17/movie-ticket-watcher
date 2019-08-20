package subscriptions

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

// MoveToArchive moves a subscription to the archives table
func MoveToArchive(dbClient dynamodbiface.ClientAPI, s Subscription) error {
	writables := []db.Writable{s}
	if err := db.Write(dbClient, writables, ArchiveTableName); err != nil {
		return fmt.Errorf("error while moving subscription to archive: %+v", err)
	}
	return Delete(dbClient, s)
}

// Delete simply deletes the given subscription from the DB
func Delete(dbClient dynamodbiface.ClientAPI, s Subscription) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"WebPushSubscription": {S: aws.String(s.WebPushSubscription)},
			"CreatedAt":           {N: aws.String(s.CreatedAt.UnixStr())},
		},
		TableName: aws.String(TableName),
	}
	req := dbClient.DeleteItemRequest(input)
	if _, err := req.Send(context.TODO()); err != nil {
		return err
	}
	return nil
}
