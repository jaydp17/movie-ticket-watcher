package db

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"reflect"
	"strconv"
	"time"
)

const timeFormat = "2006-01-02"

type YYYYMMDDTime struct {
	time.Time
}

func (t YYYYMMDDTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.ToYYYYMMDD())), nil
}

func (t *YYYYMMDDTime) UnmarshalJSON(p []byte) error {
	input, unquoteErr := strconv.Unquote(string(p))
	if unquoteErr != nil {
		return &json.UnmarshalTypeError{
			Value: input,
			Type:  reflect.TypeOf(t),
		}
	}
	parsedTime, err := time.Parse(timeFormat, input)
	if err != nil {
		return &json.UnmarshalTypeError{
			Value: input,
			Type:  reflect.TypeOf(t),
		}
	}
	t.Time = parsedTime
	return nil
}

// MarshalDynamoDBAttributeValue implements the Marshaler interface so that
// the YYYYMMDDTime can be marshaled from to a DynamoDB AttributeValue number
// value encoded in the number of seconds since January 1, 1970 UTC.
func (t YYYYMMDDTime) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	s := t.ToYYYYMMDD()
	av.S = &s
	return nil
}

// UnmarshalDynamoDBAttributeValue implements the Unmarshaler interface so that
// the YYYYMMDDTime can be unmarshaled from a DynamoDB AttributeValue number representing
// the number of seconds since January 1, 1970 UTC.
//
// If an error parsing the AttributeValue number occurs UnmarshalError will be
// returned.
func (t *YYYYMMDDTime) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	parsedTime, err := time.Parse(timeFormat, aws.StringValue(av.S))
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

func (t YYYYMMDDTime) ToYYYYMMDD() string {
	return t.Format(timeFormat)
}
