package cinemas

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
)

type CinemaResult struct {
	Cinema Cinema
	Err    error
}

func FindByID(ctx context.Context, dbClient dynamodbiface.ClientAPI, ID string) <-chan CinemaResult {
	outputCh := make(chan CinemaResult)
	if len(ID) == 0 {
		outputCh <- CinemaResult{Cinema{}, httperror.New(400, "cinemaID can't be empty")}
		close(outputCh)
		return outputCh
	}
	go func(outputCh chan<- CinemaResult) {
		defer close(outputCh)
		input := &dynamodb.GetItemInput{
			Key: map[string]dynamodb.AttributeValue{
				"ID": {S: aws.String(ID)},
			},
			TableName: aws.String(TableName),
		}
		req := dbClient.GetItemRequest(input)
		res, err := req.Send(ctx)
		if err != nil {
			outputCh <- CinemaResult{Cinema{}, err}
			return
		}
		if res.Item == nil {
			outputCh <- CinemaResult{Cinema{}, httperror.New(404, "cinema not found")}
			return
		}
		var cinema Cinema
		if err := dynamodbattribute.UnmarshalMap(res.Item, &cinema); err != nil {
			outputCh <- CinemaResult{Cinema{}, err}
			return
		}
		outputCh <- CinemaResult{cinema, nil}
	}(outputCh)
	return outputCh
}
