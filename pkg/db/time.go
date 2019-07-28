package db

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"strconv"
	"time"
)

// UnixTime was created because we wanted to control how the Marshalling & Unmarshalling happens of the time.Time type
// both is respect to the JSON marshalling & DynamoDB marshalling
type UnixTime struct {
	time.Time
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(t.UnixStr()), nil
}

func (t *UnixTime) UnmarshalJSON(p []byte) error {
	unixTime, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(unixTime, 0)
	return nil
}

// MarshalDynamoDBAttributeValue implements the Marshaler interface so that
// the UnixTime can be marshaled from to a DynamoDB AttributeValue number
// value encoded in the number of seconds since January 1, 1970 UTC.
func (t UnixTime) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	s := t.UnixStr()
	av.N = &s
	return nil
}

// UnmarshalDynamoDBAttributeValue implements the Unmarshaler interface so that
// the UnixTime can be unmarshaled from a DynamoDB AttributeValue number representing
// the number of seconds since January 1, 1970 UTC.
//
// If an error parsing the AttributeValue number occurs UnmarshalError will be
// returned.
func (t *UnixTime) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	unixTime, err := strconv.ParseInt(aws.StringValue(av.N), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(unixTime, 0)
	return nil
}

func (t UnixTime) UnixStr() string {
	unixTime := strconv.FormatInt(t.Unix(), 10)
	return unixTime
}
