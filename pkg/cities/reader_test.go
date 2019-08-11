package cities

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"testing"
)

type fakeDynamoDB struct {
	dynamodbiface.ClientAPI
	pages [][]City
	index int
	err   error
}

func (fd *fakeDynamoDB) ScanRequest(input *dynamodb.ScanInput) dynamodb.ScanRequest {
	req := dynamodb.ScanRequest{
		Input: input,
		Copy: func(v *dynamodb.ScanInput) dynamodb.ScanRequest {
			r := fd.ClientAPI.ScanRequest(v)
			r.Handlers.Clear()
			r.Handlers.Send.PushBack(func(r *aws.Request) {
				page := fd.pages[fd.index]
				output := &dynamodb.ScanOutput{Items: []map[string]dynamodb.AttributeValue{}}
				for _, city := range page {
					cityAttributeValues, _ := city.DynamoAttributeValues()
					output.Items = append(output.Items, cityAttributeValues)
				}
				if fd.index < len(fd.pages)-1 {
					output.LastEvaluatedKey = map[string]dynamodb.AttributeValue{
						"key": {S: aws.String("marker")},
					}
				}
				r.Data = output
				fd.index++
			})
			return r
		},
	}
	return req
}

func TestAll(t *testing.T) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		t.Fatalf("aws config load failed: %v", err)
	}
	fd := &fakeDynamoDB{
		ClientAPI: dynamodb.New(cfg),
		pages: [][]City{
			{
				{
					ID:           "cityID1",
					Name:         "name1",
					BookmyshowID: "bms1",
					PaytmID:      "ptm1",
					IsTopCity:    true,
				},
				{
					ID:           "cityID2",
					Name:         "name2",
					BookmyshowID: "bms2",
					PaytmID:      "ptm2",
					IsTopCity:    false,
				},
			},
			{
				{
					ID:           "cityID3",
					Name:         "name3",
					BookmyshowID: "bms3",
					PaytmID:      "ptm3",
					IsTopCity:    true,
				},
			},
		},
	}

	cityChan := All(fd)
	cityResults := make([]City, 0)
	for city := range cityChan {
		cityResults = append(cityResults, city)
	}
	if len(cityResults) != 3 {
		t.Fatalf("we didn't paginate over all the results")
	}
}
