package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"log"
	"math"
	"sync"
)

type Writable interface {
	DynamoAttributeValues() (map[string]dynamodb.AttributeValue, error)
}

func Write(dbClient dynamodbiface.ClientAPI, writables []Writable, tableName string) error {
	writeInputs := dynamoBatchWriteInputs(writables, tableName)
	errorsCh := make(chan error)
	wg := sync.WaitGroup{}
	for _, inputCopy := range writeInputs {
		wg.Add(1)
		go func(input *dynamodb.BatchWriteItemInput, errCh chan<- error) {
			defer wg.Done()
			if err := BatchWrite(dbClient, input); err != nil {
				errCh <- err
			}
		}(inputCopy, errorsCh)
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

func dynamoBatchWriteInputs(writables []Writable, tableName string) []*dynamodb.BatchWriteItemInput {
	writeRequests := make([]dynamodb.WriteRequest, 0, len(writables))
	for _, city := range writables {
		writeRequests = append(writeRequests, dynamodb.WriteRequest{PutRequest: dynamoPutReq(city)})
	}

	totalRequests := len(writeRequests)
	batches := int(math.Ceil(float64(len(writeRequests)) / MaxBatchWriteItems))
	writeInputs := make([]*dynamodb.BatchWriteItemInput, 0, batches)
	for i := 0; i < batches; i++ {
		from := i * MaxBatchWriteItems
		to := int(math.Min(float64(totalRequests), float64((i+1)*MaxBatchWriteItems)))
		chunkWriteRequests := make([]dynamodb.WriteRequest, 0, MaxBatchWriteItems)
		chunkWriteRequests = append(chunkWriteRequests, writeRequests[from:to]...)
		chunkRequestItems := map[string][]dynamodb.WriteRequest{}
		chunkRequestItems[tableName] = chunkWriteRequests
		chunkWriteInput := dynamodb.BatchWriteItemInput{RequestItems: chunkRequestItems}
		writeInputs = append(writeInputs, &chunkWriteInput)
	}
	return writeInputs
}

func dynamoPutReq(w Writable) *dynamodb.PutRequest {
	putReq, err := w.DynamoAttributeValues()
	if err != nil {
		log.Fatalln(err)
	}
	return &dynamodb.PutRequest{Item: putReq}
}
