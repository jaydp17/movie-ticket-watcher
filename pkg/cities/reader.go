package cities

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

func All() <-chan City {
	scanInput := &dynamodb.ScanInput{
		TableName:                aws.String(TableName),
		ExpressionAttributeNames: map[string]string{"#NM": "Name", "#TC": "IsTopCity"},
		ProjectionExpression:     aws.String("ID, #NM, #TC"),
	}
	req := db.Client.ScanRequest(scanInput)
	paginator := dynamodb.NewScanPaginator(req)

	pages := make(chan []map[string]dynamodb.AttributeValue)
	cities := make(chan City)

	go func() {
		for items := range pages {
			for _, item := range items {
				var city City
				if err := dynamodbattribute.UnmarshalMap(item, &city); err != nil {
					fmt.Println("Got error unmarshalling:")
					fmt.Println(err.Error())
					return
				}
				cities <- city
			}
		}
		close(cities)
	}()

	for paginator.Next(context.TODO()) {
		page := paginator.CurrentPage()
		pages <- page.Items
	}
	close(pages)

	if err := paginator.Err(); err != nil {
		fmt.Println("error in paginator")
		fmt.Println(err)
	}

	return cities
}
