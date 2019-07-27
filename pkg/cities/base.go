package cities

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/config"
)

type City struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	BookmyshowID string `json:"bookmyshowID,omitempty"`
	PaytmID      string `json:"paytmID,omitempty"`
	IsTopCity    bool   `json:"isTopCity"`
}

var TableName = config.FullTableName("cities")

// HasAllProviderIDs checks if a cities object has all the provider IDs
func (c City) HasAllProviderIDs() bool {
	return len(c.BookmyshowID) > 0 && len(c.PaytmID) > 0
}

func (c City) DynamoAttributeValues() (map[string]dynamodb.AttributeValue, error) {
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
