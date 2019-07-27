package moviecitylink

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/config"
)

type MovieCityLink struct {
	CityID   string   `json:"cityID"`
	MovieIDs []string `json:"movieIDs"`
}

var TableName = config.FullTableName("movie-city-link")

func (link MovieCityLink) DynamoAttributeValues() (map[string]dynamodb.AttributeValue, error) {
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
