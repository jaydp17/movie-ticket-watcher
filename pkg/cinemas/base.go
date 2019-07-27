package cinemas

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/config"
)

type Cinema struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Provider     string  `json:"provider"`
	CityID       string  `json:"cityID"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Address      string  `json:"address"`
	BookmyshowID string  `json:"bookmyshowID,omitempty"`
	PaytmID      string  `json:"paytmID,omitempty"`
}

var TableName = config.FullTableName("cinemas")

func (c Cinema) HasAllProviderIDs() bool {
	return len(c.BookmyshowID) > 0 && len(c.PaytmID) > 0
}

func (c Cinema) DynamoAttributeValues() (map[string]dynamodb.AttributeValue, error) {
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
