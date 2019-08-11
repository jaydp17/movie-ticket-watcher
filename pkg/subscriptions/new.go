package subscriptions

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
	"time"
)

// New validates all the city/movie/cinema IDs & creates a new Subscription object from those IDs
func New(dbClient dynamodbiface.ClientAPI, cityID, movieID, cinemaID, webPushSubscription string, date time.Time) (Subscription, error) {
	if date.IsZero() {
		return Subscription{}, httperror.New(400, "date can't be zero")
	}
	if time.Now().After(date) {
		return Subscription{}, httperror.New(400, "date can't be in the past")
	}
	city, movie, cinema, err := getMovieCityAndCinema(dbClient, cityID, movieID, cinemaID)
	if err != nil {
		return Subscription{}, err
	}
	return Subscription{
		WebPushSubscription: webPushSubscription,
		CreatedAt:           db.UnixTime{Time: time.Now()},
		CityID:              city.ID,
		MovieID:             movie.ID,
		CinemaID:            cinema.ID,
		ScreeningDate:       db.YYYYMMDDTime{Time: date},
	}, nil
}
