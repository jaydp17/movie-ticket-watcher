package movies

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
)

type MovieResult struct {
	Movie Movie
	Err   error
}

func FindByID(ctx context.Context, dbClient dynamodbiface.ClientAPI, ID string) <-chan MovieResult {
	outputCh := make(chan MovieResult)
	go func(outputCh chan<- MovieResult) {
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
			outputCh <- MovieResult{Movie{}, err}
			return
		}

		if res.Item == nil {
			outputCh <- MovieResult{Movie{}, httperror.New(404, "movie not found")}
			return
		}
		var movie Movie
		if err := dynamodbattribute.UnmarshalMap(res.Item, &movie); err != nil {
			outputCh <- MovieResult{Movie{}, err}
			return
		}

		outputCh <- MovieResult{movie, nil}
	}(outputCh)
	return outputCh
}
