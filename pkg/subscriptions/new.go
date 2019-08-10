package subscriptions

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
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

func getMovieCityAndCinema(dbClient dynamodbiface.ClientAPI, cityID, movieID, cinemaID string) (cities.City, movies.Movie, cinemas.Cinema, error) {
	if len(cityID) == 0 {
		return cities.City{}, movies.Movie{}, cinemas.Cinema{}, httperror.New(400, "cityID can't be empty")
	}
	if len(movieID) == 0 {
		return cities.City{}, movies.Movie{}, cinemas.Cinema{}, httperror.New(400, "movieID can't be empty")
	}
	if len(cinemaID) == 0 {
		return cities.City{}, movies.Movie{}, cinemas.Cinema{}, httperror.New(400, "cinemaID can't be empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cityOutput := cities.FindByID(ctx, dbClient, cityID)
	movieOutput := movies.FindByID(ctx, dbClient, movieID)
	cinemaOutput := cinemas.FindByID(ctx, dbClient, cinemaID)

	var city cities.City
	var movie movies.Movie
	var cinema cinemas.Cinema
	for cityOutput != nil || movieOutput != nil || cinemaOutput != nil {
		//noinspection GoNilness
		select {
		case result := <-cityOutput:
			if result.Err != nil {
				cancel()
				return cities.City{}, movies.Movie{}, cinemas.Cinema{}, result.Err
			}
			city = result.City
			cityOutput = nil
		case result := <-movieOutput:
			if result.Err != nil {
				cancel()
				return cities.City{}, movies.Movie{}, cinemas.Cinema{}, result.Err
			}
			movie = result.Movie
			movieOutput = nil
		case result := <-cinemaOutput:
			if result.Err != nil {
				cancel()
				return cities.City{}, movies.Movie{}, cinemas.Cinema{}, result.Err
			}
			cinema = result.Cinema
			cinemaOutput = nil
		}
	}
	return city, movie, cinema, nil
}
