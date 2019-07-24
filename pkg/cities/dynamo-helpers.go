package cities

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"log"
	"math"
)

func dynamoBatchWriteInputs(cities []City) []*dynamodb.BatchWriteItemInput {
	writeRequests := make([]dynamodb.WriteRequest, 0, len(cities))
	for _, city := range cities {
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

func dynamoPutReq(c City) *dynamodb.PutRequest {
	putReq, err := dynamoAttributeValues(c)
	if err != nil {
		log.Fatalln(err)
	}
	return &dynamodb.PutRequest{Item: putReq}
}

func dynamoAttributeValues(c City) (map[string]dynamodb.AttributeValue, error) {
	if len(c.ID) == 0 {
		return nil, fmt.Errorf("ID is empty: %+v", c)
	}
	if len(c.Name) == 0 {
		return nil, fmt.Errorf("name is empty: %+v", c)
	}
	if len(c.BookmyshowID) == 0 {
		return nil, fmt.Errorf("BookmyshowID is empty: %+v", c)
	}
	if len(c.PaytmID) == 0 {
		return nil, fmt.Errorf("PaytmID is empty: %+v", c)
	}
	item := map[string]dynamodb.AttributeValue{
		"ID":           {S: aws.String(c.ID)},
		"Name":         {S: aws.String(c.Name)},
		"BookmyshowID": {S: aws.String(c.BookmyshowID)},
		"PaytmID":      {S: aws.String(c.PaytmID)},
		"IsTopCity":    {BOOL: aws.Bool(c.IsTopCity)},
	}
	return item, nil
}
