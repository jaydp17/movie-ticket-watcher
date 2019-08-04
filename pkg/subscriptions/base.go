package subscriptions

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/config"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

type Subscription struct {
	ID        string      `json:"id"`
	CityID    string      `json:"cityID"`
	MovieID   string      `json:"movieID"`
	CinemaID  string      `json:"cinemaID"`
	CreatedAt db.UnixTime `json:"createdAt"`
}

var TableName = config.FullTableName("subscriptions")

func (s Subscription) DynamoAttributeValues() (map[string]dynamodb.AttributeValue, error) {
	if len(s.ID) == 0 {
		return nil, fmt.Errorf("ID is empty: %+v", s)
	}
	if len(s.CityID) == 0 {
		return nil, fmt.Errorf("cityID is empty: %+v", s)
	}
	if len(s.MovieID) == 0 {
		return nil, fmt.Errorf("movieID is empty: %+v", s)
	}
	if len(s.CinemaID) == 0 {
		return nil, fmt.Errorf("cinemaID is empty: %+v", s)
	}

	item := map[string]dynamodb.AttributeValue{
		"ID":        {S: aws.String(s.ID)},
		"CityID":    {S: aws.String(s.CityID)},
		"MovieID":   {S: aws.String(s.MovieID)},
		"CinemaID":  {S: aws.String(s.CinemaID)},
		"CreatedAt": {N: aws.String(s.CreatedAt.UnixStr())},
	}
	return item, nil
}