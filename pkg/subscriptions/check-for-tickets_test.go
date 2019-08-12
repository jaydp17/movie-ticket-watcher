package subscriptions

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type mockProviderForCheckForAvailableTickets struct {
	result providers.VenueCodesResult
}

func (p mockProviderForCheckForAvailableTickets) FetchAvailableVenueCodes(providerCityID, providerMovieID string, date db.YYYYMMDDTime) <-chan providers.VenueCodesResult {
	ch := make(chan providers.VenueCodesResult)
	go func() {
		ch <- p.result
		close(ch)
	}()
	return ch
}

func TestCheckForAvailableTickets(t *testing.T) {
	type VenueCodesResult = providers.VenueCodesResult
	dbClient := newFakeGetCityMovieCinemaDB()

	t.Run("single available subscription", func(t *testing.T) {
		s := Subscription{
			CityID:              "BANG",
			MovieID:             "endgame",
			CinemaID:            "pvrxm",
			WebPushSubscription: "abcd",
			ScreeningDate:       db.YYYYMMDDTime{Time: time.Now()},
		}
		cinema := dbClient.cinema

		bms := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Data: []string{cinema.BookmyshowID}}}
		ptm := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Data: []string{cinema.PaytmID}}}

		availableSubscriptionsCh := CheckForAvailableTickets(dbClient, bms, ptm, []Subscription{s})
		result := make([]Subscription, 0)
		for sub := range availableSubscriptionsCh {
			result = append(result, sub)
		}
		assert.Equal(t, result, []Subscription{s})
	})

	t.Run("single available subscription from BMS only", func(t *testing.T) {
		s := Subscription{
			CityID:              "BANG",
			MovieID:             "endgame",
			CinemaID:            "pvrxm",
			WebPushSubscription: "abcd",
			ScreeningDate:       db.YYYYMMDDTime{Time: time.Now()},
		}
		cinema := dbClient.cinema

		bms := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Data: []string{cinema.BookmyshowID}}}
		ptm := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Err: fmt.Errorf("request failed")}}

		availableSubscriptionsCh := CheckForAvailableTickets(dbClient, bms, ptm, []Subscription{s})
		result := make([]Subscription, 0)
		for sub := range availableSubscriptionsCh {
			result = append(result, sub)
		}
		assert.Equal(t, result, []Subscription{s})
	})

	t.Run("single available subscription from PTM only", func(t *testing.T) {
		s := Subscription{
			CityID:              "BANG",
			MovieID:             "endgame",
			CinemaID:            "pvrxm",
			WebPushSubscription: "abcd",
			ScreeningDate:       db.YYYYMMDDTime{Time: time.Now()},
		}
		cinema := dbClient.cinema

		bms := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Err: fmt.Errorf("request failed")}}
		ptm := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Data: []string{cinema.PaytmID}}}

		availableSubscriptionsCh := CheckForAvailableTickets(dbClient, bms, ptm, []Subscription{s})
		result := make([]Subscription, 0)
		for sub := range availableSubscriptionsCh {
			result = append(result, sub)
		}
		assert.Equal(t, result, []Subscription{s})
	})

	t.Run("single un-available subscription", func(t *testing.T) {
		s := Subscription{
			CityID:              "BANG",
			MovieID:             "endgame",
			CinemaID:            "pvrxm",
			WebPushSubscription: "abcd",
			ScreeningDate:       db.YYYYMMDDTime{Time: time.Now()},
		}

		bms := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Err: fmt.Errorf("request failed")}}
		ptm := mockProviderForCheckForAvailableTickets{result: VenueCodesResult{Err: fmt.Errorf("request failed")}}

		availableSubscriptionsCh := CheckForAvailableTickets(dbClient, bms, ptm, []Subscription{s})
		result := make([]Subscription, 0)
		for sub := range availableSubscriptionsCh {
			result = append(result, sub)
		}
		assert.Equal(t, result, []Subscription{})
	})
}

func TestGroupSimilarSubscriptions(t *testing.T) {
	date := db.YYYYMMDDTime{Time: time.Now()}
	input := []Subscription{
		{
			CityID:        "a",
			MovieID:       "b",
			ScreeningDate: date,
		},
		{
			CityID:        "a",
			MovieID:       "b",
			ScreeningDate: date,
		},
		{
			CityID:        "b",
			MovieID:       "c",
			ScreeningDate: date,
		},
	}
	expectedOutput := []groupSubscriptions{
		{
			subscriptions: []Subscription{
				{
					CityID:        "a",
					MovieID:       "b",
					ScreeningDate: date,
				},
				{
					CityID:        "a",
					MovieID:       "b",
					ScreeningDate: date,
				},
			},
		},
		{
			subscriptions: []Subscription{
				{
					CityID:        "b",
					MovieID:       "c",
					ScreeningDate: date,
				},
			},
		},
	}

	result := groupSimilarSubscriptions(input)
	if len(result) != len(expectedOutput) {
		t.Errorf("expectedOutput(%d groups) & result(%d groups) don't have the same number of groups", len(expectedOutput), len(result))
	}

	for i := 0; i < len(result); i++ {
		resultGroup := result[i].subscriptions
		expectedGroup := expectedOutput[i].subscriptions
		if len(resultGroup) != len(expectedGroup) {
			t.Errorf("subscriptions in the groups don't match with the expected output")
		}
		for j := 0; j < len(resultGroup); j++ {
			if resultGroup[j].CityID != expectedGroup[j].CityID || resultGroup[j].MovieID != expectedGroup[j].MovieID || resultGroup[j].ScreeningDate.ToYYYYMMDD() != expectedGroup[j].ScreeningDate.ToYYYYMMDD() {
				t.Errorf("subscriptions in the groups don't match")
			}
		}
	}
}

func newFakeGetCityMovieCinemaDB() *fakeGetCityMovieCinemaDB {
	city := cities.City{
		ID:           "BANG",
		Name:         "Bengaluru",
		BookmyshowID: "BANG",
		PaytmID:      "bengaluru",
	}
	movie := movies.Movie{
		ID:           "endgame",
		GroupID:      "endgame",
		Title:        "Avengers-endgame",
		ScreenFormat: "3d",
		Language:     "English",
		ImageURL:     "",
		BookmyshowID: "MX012345",
		PaytmID:      "QF012345",
		TTL:          db.UnixTime{Time: time.Now()},
	}
	cinema := cinemas.Cinema{
		ID:           "cinemaID",
		Name:         "PVR bellendur",
		Provider:     "PVR",
		CityID:       city.ID,
		Latitude:     1,
		Longitude:    2,
		BookmyshowID: "PVRMX",
		PaytmID:      "123",
	}
	dbClient := &fakeGetCityMovieCinemaDB{
		city:   city,
		movie:  movie,
		cinema: cinema,
	}
	return dbClient
}

type fakeGetCityMovieCinemaDB struct {
	dynamodbiface.ClientAPI
	city   cities.City
	movie  movies.Movie
	cinema cinemas.Cinema
}

func (fd *fakeGetCityMovieCinemaDB) GetItemRequest(input *dynamodb.GetItemInput) dynamodb.GetItemRequest {
	var item map[string]dynamodb.AttributeValue
	var dynamoMarshalErr error
	if *input.TableName == cities.TableName {
		item, dynamoMarshalErr = fd.city.DynamoAttributeValues()
	} else if *input.TableName == movies.TableName {
		item, dynamoMarshalErr = fd.movie.DynamoAttributeValues()
	} else if *input.TableName == cinemas.TableName {
		item, dynamoMarshalErr = fd.cinema.DynamoAttributeValues()
	}
	if dynamoMarshalErr != nil {
		panic(dynamoMarshalErr)
	}
	output := &dynamodb.GetItemOutput{Item: item}
	req := dynamodb.GetItemRequest{
		Request: &aws.Request{
			Data:        output,
			Error:       dynamoMarshalErr,
			HTTPRequest: &http.Request{},
		},
	}
	return req
}
