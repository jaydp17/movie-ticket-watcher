package dao

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"log"
	"math"
	"sync"
)

type City struct {
	ID           string
	Name         string
	BookmyshowID string
	PaytmID      string
	IsTopCity    bool
}

type Cities []City

var CityTableName = db.GetFullTableName("cities")

// HasAllProviderIDs checks if a city object has all the provider IDs
func (c *City) HasAllProviderIDs() bool {
	return len(c.BookmyshowID) > 0 && len(c.PaytmID) > 0
}

// Write allows us to write a list of cities in the DB
// It divides all the Items (i.e. documents) in batches of 25 (i.e. the max limit of items DynamoDB accepts per batch request)
// and then concurrently inserts those batches in the DB
func (cities Cities) Write() error {
	writeInputs := cities.dynamoBatchWriteInputs()
	errorsCh := make(chan error)
	wg := sync.WaitGroup{}
	for _, input := range writeInputs {
		wg.Add(1)
		go func(errCh chan<- error) {
			defer wg.Done()
			if err := db.BatchWrite(input); err != nil {
				errCh <- err
			}
		}(errorsCh)
	}

	go func() {
		wg.Wait()
		close(errorsCh)
	}()

	// collect all errors from the channel
	var errors []error
	for err := range errorsCh {
		errors = append(errors, err)
	}
	// concat all errors into one
	if len(errors) > 0 {
		errorsStr := fmt.Sprintf("failed to write because of the following errors:\n---------")
		for _, err := range errors {
			errorsStr = fmt.Sprintf("%s\n%+v", errorsStr, err)
		}
		return fmt.Errorf(errorsStr)
	}
	return nil
}

func (cities Cities) dynamoBatchWriteInputs() []*dynamodb.BatchWriteItemInput {
	writeRequests := make([]dynamodb.WriteRequest, 0, len(cities))
	for _, city := range cities {
		writeRequests = append(writeRequests, dynamodb.WriteRequest{PutRequest: city.dynamoPutReq()})
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
		chunkRequestItems[CityTableName] = chunkWriteRequests
		chunkWriteInput := dynamodb.BatchWriteItemInput{RequestItems: chunkRequestItems}
		writeInputs = append(writeInputs, &chunkWriteInput)
	}
	return writeInputs
}

func (c *City) dynamoPutReq() *dynamodb.PutRequest {
	putReq, err := c.dynamoAttributeValues()
	if err != nil {
		log.Fatalln(err)
	}
	return &dynamodb.PutRequest{Item: putReq}
}

func (c *City) dynamoAttributeValues() (map[string]dynamodb.AttributeValue, error) {
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
