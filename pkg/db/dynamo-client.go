package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var Client *dynamodb.Client

func init() {
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
	Client = dynamodb.New(awsConfig)
}

// ListTables allows you to list all the tables in this region
func ListTables() {
	req := Client.ListTablesRequest(&dynamodb.ListTablesInput{})

	// Send the request, and get the response or error back
	resp, err := req.Send(context.Background())
	if err != nil {
		panic("failed to describe table, " + err.Error())
	}

	fmt.Println("Response", resp)
}
