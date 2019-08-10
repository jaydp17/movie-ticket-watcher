package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/moviecitylink"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
)

func main() {
	dbClient := db.NewClient()

	createCitiesTable(dbClient)
	createMoviesTable(dbClient)
	createCinemasTable(dbClient)
	createMovieCityLink(dbClient)
	createSubscriptionsTable(dbClient)
}

func createCitiesTable(dbClient dynamodbiface.ClientAPI) {
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
	req := dbClient.CreateTableRequest(input)
	sendReq(&req)
}

func createMoviesTable(dbClient dynamodbiface.ClientAPI) {
	tableName := movies.TableName
	createTableInput := &dynamodb.CreateTableInput{
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

	req := dbClient.CreateTableRequest(createTableInput)
	sendReq(&req)

	describeTableInput := dynamodb.DescribeTableInput{TableName: aws.String(tableName)}
	if err := dbClient.WaitUntilTableExists(context.TODO(), &describeTableInput); err != nil {
		fmt.Printf("error waitng for table to be created: %v", err)
		return
	}

	ttlInput := dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String("TTL"),
			Enabled:       aws.Bool(true),
		},
	}
	ttlReq := dbClient.UpdateTimeToLiveRequest(&ttlInput)
	_, err := ttlReq.Send(context.TODO())
	if err != nil {
		fmt.Printf("error creating TTL attribute: %v", err)
		return
	}
}

func createCinemasTable(dbClient dynamodbiface.ClientAPI) {
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
	req := dbClient.CreateTableRequest(input)
	sendReq(&req)
}

func createMovieCityLink(dbClient dynamodbiface.ClientAPI) {
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
	req := dbClient.CreateTableRequest(input)
	sendReq(&req)
}

func createSubscriptionsTable(dbClient dynamodbiface.ClientAPI) {
	tableName := subscriptions.TableName
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("WebPushSubscription"),
				AttributeType: dynamodb.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("CreatedAt"),
				AttributeType: dynamodb.ScalarAttributeTypeN,
			},
		},
		KeySchema: []dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("WebPushSubscription"),
				KeyType:       dynamodb.KeyTypeHash,
			},
			{
				AttributeName: aws.String("CreatedAt"),
				KeyType:       dynamodb.KeyTypeRange,
			},
		},
		BillingMode: dynamodb.BillingModePayPerRequest,
	}
	req := dbClient.CreateTableRequest(input)
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
