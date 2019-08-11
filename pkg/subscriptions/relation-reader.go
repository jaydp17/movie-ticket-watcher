package subscriptions

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
)

func (s Subscription) GetMovieCityAndCinema(dbClient dynamodbiface.ClientAPI) (cities.City, movies.Movie, cinemas.Cinema, error) {
	return getMovieCityAndCinema(dbClient, s.CityID, s.MovieID, s.CinemaID)
}

func getMovieCityAndCinema(dbClient dynamodbiface.ClientAPI, cityID, movieID, cinemaID string) (cities.City, movies.Movie, cinemas.Cinema, error) {
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
