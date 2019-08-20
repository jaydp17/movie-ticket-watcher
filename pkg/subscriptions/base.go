package subscriptions

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/config"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

type Subscription struct {
	CityID              string          `json:"cityID"`
	MovieID             string          `json:"movieID"`
	CinemaID            string          `json:"cinemaID"`
	WebPushSubscription string          `json:"webPushSubscription"`
	ScreeningDate       db.YYYYMMDDTime `json:"screeningDate"`
	CreatedAt           db.UnixTime     `json:"createdAt"`
}

var TableName = config.FullTableName("subscriptions")
var ArchiveTableName = config.FullTableName("subscriptions-archive")

// IsSimilar just checks for CityID, MovieID & ScreeningDate
func (s Subscription) IsSimilar(s2 Subscription) bool {
	return s.CityID == s2.CityID && s.MovieID == s2.MovieID && s.ScreeningDate.ToYYYYMMDD() == s2.ScreeningDate.ToYYYYMMDD()
}

func (s Subscription) DynamoAttributeValues() (map[string]dynamodb.AttributeValue, error) {
	if len(s.WebPushSubscription) == 0 {
		return nil, fmt.Errorf("webPushSubscription is empty: %+v", s)
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
	if s.ScreeningDate.IsZero() {
		return nil, fmt.Errorf("screeningDate is zero: %+v", s)
	}

	item := map[string]dynamodb.AttributeValue{
		"WebPushSubscription": {S: aws.String(s.WebPushSubscription)},
		"CreatedAt":           {N: aws.String(s.CreatedAt.UnixStr())},
		"CityID":              {S: aws.String(s.CityID)},
		"MovieID":             {S: aws.String(s.MovieID)},
		"CinemaID":            {S: aws.String(s.CinemaID)},
		"ScreeningDate":       {S: aws.String(s.ScreeningDate.ToYYYYMMDD())},
	}
	return item, nil
}
