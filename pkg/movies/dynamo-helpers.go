package movies

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"log"
	"math"
)

func dynamoBatchWriteInputs(movies []Movie) []*dynamodb.BatchWriteItemInput {
	writeRequests := make([]dynamodb.WriteRequest, 0, len(movies))
	for _, city := range movies {
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

func dynamoPutReq(m Movie) *dynamodb.PutRequest {
	putReq, err := dynamoAttributeValues(m)
	if err != nil {
		log.Fatalln(err)
	}
	return &dynamodb.PutRequest{Item: putReq}
}

func dynamoAttributeValues(m Movie) (map[string]dynamodb.AttributeValue, error) {
	if len(m.ID) == 0 {
		return nil, fmt.Errorf("ID is empty: %+v", m)
	}
	if len(m.GroupID) == 0 {
		return nil, fmt.Errorf("groupID is empty: %+v", m)
	}
	if len(m.Title) == 0 {
		return nil, fmt.Errorf("title is empty: %+v", m)
	}
	if len(m.ScreenFormat) == 0 {
		return nil, fmt.Errorf("screenFormat is empty: %+v", m)
	}
	if len(m.Language) == 0 {
		return nil, fmt.Errorf("language is empty: %+v", m)
	}
	if len(m.BookmyshowID) == 0 && len(m.PaytmID) == 0 {
		return nil, fmt.Errorf("BookmyshowID & PaytmID both can't be empty: %+v", m)
	}
	item := map[string]dynamodb.AttributeValue{
		"ID":           {S: aws.String(m.ID)},
		"GroupID":      {S: aws.String(m.GroupID)},
		"Title":        {S: aws.String(m.Title)},
		"ScreenFormat": {S: aws.String(m.ScreenFormat)},
		"Language":     {S: aws.String(m.Language)},
	}
	if len(m.ImageURL) > 0 {
		item["ImageURL"] = dynamodb.AttributeValue{S: aws.String(m.ImageURL)}
	}
	if len(m.BookmyshowID) > 0 {
		item["BookmyshowID"] = dynamodb.AttributeValue{S: aws.String(m.BookmyshowID)}
	}
	if len(m.PaytmID) > 0 {
		item["PaytmID"] = dynamodb.AttributeValue{S: aws.String(m.PaytmID)}
	}
	return item, nil
}
