package moviecitylink

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"log"
	"math"
)

func dynamoBatchWriteInputs(links []MovieCityLink) []*dynamodb.BatchWriteItemInput {
	writeRequests := make([]dynamodb.WriteRequest, 0, len(links))
	for _, city := range links {
		writeRequests = append(writeRequests, dynamodb.WriteRequest{PutRequest: dynamoPutReq(city)})
	}

	totalRequests := len(writeRequests)
	batches := int(math.Ceil(float64(len(writeRequests)) / db.MaxBatchWriteItems))
	writeInputs := make([]*dynamodb.BatchWriteItemInput, 0, batches)
	for i := 0; i < batches; i++ {
		from := i * db.MaxBatchWriteItems
		to := int(math.Min(float64(totalRequests), float64((i+1)*db.MaxBatchWriteItems)))
		chunkWriteRequests := make([]dynamodb.WriteRequest, 0, db.MaxBatchWriteItems)
		chunkWriteRequests = append(chunkWriteRequests, writeRequests[from:to]...)
		chunkRequestItems := map[string][]dynamodb.WriteRequest{}
		chunkRequestItems[TableName] = chunkWriteRequests
		chunkWriteInput := dynamodb.BatchWriteItemInput{RequestItems: chunkRequestItems}
		writeInputs = append(writeInputs, &chunkWriteInput)
	}
	return writeInputs
}

func dynamoPutReq(link MovieCityLink) *dynamodb.PutRequest {
	putReq, err := dynamoAttributeValues(link)
	if err != nil {
		log.Fatalln(err)
	}
	return &dynamodb.PutRequest{Item: putReq}
}

func dynamoAttributeValues(link MovieCityLink) (map[string]dynamodb.AttributeValue, error) {
	if len(link.CityID) == 0 {
		return nil, fmt.Errorf("CityID is empty: %+v", link)
	}
	if len(link.MovieIDs) == 0 {
		return nil, fmt.Errorf("MovieIDs is empty: %+v", link)
	}
	item := map[string]dynamodb.AttributeValue{
		"CityID":   {S: aws.String(link.CityID)},
		"MovieIDs": {SS: link.MovieIDs},
	}
	return item, nil
}
