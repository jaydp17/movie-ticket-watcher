package movies

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaydp17/movie-ticket-watcher/pkg/config"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

type Movie struct {
	ID           string      `json:"id"`
	GroupID      string      `json:"groupID"`
	Title        string      `json:"title"` // make sure there's no additional information in the title like (3D) or the Language
	ScreenFormat string      `json:"screenFormat"`
	Language     string      `json:"language"`
	ImageURL     string      `json:"imageURL"`
	BookmyshowID string      `json:"bookmyshowID"`
	PaytmID      string      `json:"paytmID"`
	TTL          db.UnixTime `json:"ttl"`
}

var TableName = config.FullTableName("movies")

func (m Movie) DynamoAttributeValues() (map[string]dynamodb.AttributeValue, error) {
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
	if m.TTL.IsZero() {
		return nil, fmt.Errorf("TTL can't be zero: %+v", m)
	}
	item := map[string]dynamodb.AttributeValue{
		"ID":           {S: aws.String(m.ID)},
		"GroupID":      {S: aws.String(m.GroupID)},
		"Title":        {S: aws.String(m.Title)},
		"ScreenFormat": {S: aws.String(m.ScreenFormat)},
		"Language":     {S: aws.String(m.Language)},
		"TTL":          {N: aws.String(m.TTL.UnixStr())},
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
