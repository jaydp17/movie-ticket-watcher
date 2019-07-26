package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/moviecitylink"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
)

func main() {
	createCitiesTable()
	createMoviesTable()
	createCinemasTable()
	createMovieCityLink()
}

func createCitiesTable() {
	tableName := cities.TableName
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{

				AttributeName: aws.String("ID"),
				KeyType:       dynamodb.KeyTypeHash,
			},
		},
		BillingMode: dynamodb.BillingModePayPerRequest,
	}
	req := db.Client.CreateTableRequest(input)
	sendReq(&req)
}

func createMoviesTable() {
	tableName := movies.TableName
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{

				AttributeName: aws.String("ID"),
				KeyType:       dynamodb.KeyTypeHash,
			},
		},
		BillingMode: dynamodb.BillingModePayPerRequest,
	}

	// CityID Global Secondary Index
	//cityIndex := dynamodb.GlobalSecondaryIndex{
	//	IndexName: aws.String("city-id-index"),
	//	KeySchema: []dynamodb.KeySchemaElement{
	//		{
	//
	//			AttributeName: aws.String("cityID"),
	//			KeyType:       dynamodb.KeyTypeHash,
	//		},
	//	},
	//	Projection: &dynamodb.Projection{
	//		ProjectionType: dynamodb.ProjectionTypeAll,
	//	},
	//}
	//input.GlobalSecondaryIndexes = append(input.GlobalSecondaryIndexes, cityIndex)

	req := db.Client.CreateTableRequest(input)
	sendReq(&req)
}

func createCinemasTable() {
	tableName := cinemas.TableName
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{

				AttributeName: aws.String("ID"),
				KeyType:       dynamodb.KeyTypeHash,
			},
		},
		BillingMode: dynamodb.BillingModePayPerRequest,
	}
	req := db.Client.CreateTableRequest(input)
	sendReq(&req)
}

func createMovieCityLink() {
	tableName := moviecitylink.TableName
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("CityID"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{

				AttributeName: aws.String("CityID"),
				KeyType:       dynamodb.KeyTypeHash,
			},
		},
		BillingMode: dynamodb.BillingModePayPerRequest,
	}
	req := db.Client.CreateTableRequest(input)
	sendReq(&req)
}

// a function that takes a createTableRequest which has all the information about how to create a table
// and goes ahead & executes it
func sendReq(req *dynamodb.CreateTableRequest) {
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				fmt.Println(dynamodb.ErrCodeResourceInUseException, aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
