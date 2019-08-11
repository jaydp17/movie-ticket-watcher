package subscriptions

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
)

func All(dbClient dynamodbiface.ClientAPI) <-chan Subscription {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(TableName),
	}

	req := dbClient.ScanRequest(scanInput)
	paginator := dynamodb.NewScanPaginator(req)

	pages := make(chan []map[string]dynamodb.AttributeValue)
	subscriptions := make(chan Subscription)

	go func() {
		for page := range pages {
			for _, item := range page {
				var subscription Subscription
				if err := dynamodbattribute.UnmarshalMap(item, &subscription); err != nil {
					fmt.Printf("Got error unmarshalling: %v", err)
					return
				}
				subscriptions <- subscription
			}
		}
		close(subscriptions)
	}()

	for paginator.Next(context.TODO()) {
		page := paginator.CurrentPage()
		pages <- page.Items
	}
	close(pages)

	if err := paginator.Err(); err != nil {
		fmt.Printf("error in paginator: %v", err)
	}

	return subscriptions
}
