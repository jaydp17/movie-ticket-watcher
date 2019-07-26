package cinemas

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"log"
	"math"
)

func dynamoBatchWriteInputs(cinemas []Cinema) []*dynamodb.BatchWriteItemInput {
	writeRequests := make([]dynamodb.WriteRequest, 0, len(cinemas))
	for _, city := range cinemas {
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

func dynamoPutReq(c Cinema) *dynamodb.PutRequest {
	putReq, err := dynamoAttributeValues(c)
	if err != nil {
		log.Fatalln(err)
	}
	return &dynamodb.PutRequest{Item: putReq}
}

func dynamoAttributeValues(c Cinema) (map[string]dynamodb.AttributeValue, error) {
	if len(c.ID) == 0 {
		return nil, fmt.Errorf("ID is empty: %+v", c)
	}
	if len(c.Name) == 0 {
		return nil, fmt.Errorf("name is empty: %+v", c)
	}
	if len(c.CityID) == 0 {
		return nil, fmt.Errorf("cityID is empty: %+v", c)
	}
	if c.Latitude == 0.0 && c.Longitude == 0.0 {
		return nil, fmt.Errorf("latitude & longitude both can't be empty: %+v", c)
	}
	if len(c.BookmyshowID) == 0 && len(c.PaytmID) == 0 {
		return nil, fmt.Errorf("BookmyshowID & PaytmID both can't be empty: %+v", c)
	}
	item := map[string]dynamodb.AttributeValue{
		"ID":        {S: aws.String(c.ID)},
		"Name":      {S: aws.String(c.Name)},
		"CityID":    {S: aws.String(c.CityID)},
		"Latitude":  {N: aws.String(fmt.Sprintf("%.6f", c.Latitude))},
		"Longitude": {N: aws.String(fmt.Sprintf("%.6f", c.Longitude))},
	}
	if len(c.Provider) > 0 {
		item["Provider"] = dynamodb.AttributeValue{S: aws.String(c.Provider)}
	}
	if len(c.Address) > 0 {
		item["Address"] = dynamodb.AttributeValue{S: aws.String(c.Address)}
	}
	if len(c.BookmyshowID) > 0 {
		item["BookmyshowID"] = dynamodb.AttributeValue{S: aws.String(c.BookmyshowID)}
	}
	if len(c.PaytmID) > 0 {
		item["PaytmID"] = dynamodb.AttributeValue{S: aws.String(c.PaytmID)}
	}
	return item, nil
}
