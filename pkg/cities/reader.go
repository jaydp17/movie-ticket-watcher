package cities

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
)

type CityResult struct {
	City City
	Err  error
}

func FindByID(ctx context.Context, cityID string) <-chan CityResult {
	outputCh := make(chan CityResult)
	go func(outputCh chan<- CityResult) {
		defer close(outputCh)
		input := &dynamodb.GetItemInput{
			Key: map[string]dynamodb.AttributeValue{
				"ID": {S: aws.String(cityID)},
			},
			TableName: aws.String(TableName),
		}
		req := db.Client.GetItemRequest(input)
		result, err := req.Send(ctx)
		if err != nil {
			outputCh <- CityResult{City{}, err}
			return
		}

		if result.Item == nil {
			outputCh <- CityResult{City{}, httperror.New(404, "city not found")}
			return
		}

		var city City
		if err := dynamodbattribute.UnmarshalMap(result.Item, &city); err != nil {
			outputCh <- CityResult{City{}, err}
			return
		}
		outputCh <- CityResult{city, nil}
	}(outputCh)
	return outputCh
}

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
