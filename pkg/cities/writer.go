package cities

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"sync"
)

func Write(cities []City) error {
	writeInputs := dynamoBatchWriteInputs(cities)
	errorsCh := make(chan error)
	wg := sync.WaitGroup{}
	for _, inputCopy := range writeInputs {
		wg.Add(1)
		go func(input *dynamodb.BatchWriteItemInput, errCh chan<- error) {
			defer wg.Done()
			if err := db.BatchWrite(input); err != nil {
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
