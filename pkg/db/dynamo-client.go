package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
)

const MaxBatchWriteItems = 25

func NewClient() *dynamodb.Client {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	awsConfig, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	awsConfig.Region = endpoints.ApSouth1RegionID

	// Using the Config value, create the DynamoDB Client
	return dynamodb.New(awsConfig)
}

// ListTables allows you to list all the tables in this region
func ListTables(dbClient dynamodbiface.ClientAPI) {
	req := dbClient.ListTablesRequest(&dynamodb.ListTablesInput{})

	// Send the request, and get the response or error back
	resp, err := req.Send(context.Background())
	if err != nil {
		panic("failed to describe table, " + err.Error())
	}

	fmt.Println("Response", resp)
}

func BatchWrite(dbClient dynamodbiface.ClientAPI, input *dynamodb.BatchWriteItemInput) error {
	req := dbClient.BatchWriteItemRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
			return aerr
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return err
		}
	}
	fmt.Println(result)
	return nil
}
