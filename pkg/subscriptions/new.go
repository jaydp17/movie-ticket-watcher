package subscriptions

import (
	"context"
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
	"github.com/jaydp17/movie-ticket-watcher/pkg/utils"
	"time"
)

// New validates all the city/movie/cinema IDs & creates a new Subscription object from those IDs
func New(cityID, movieID, cinemaID string, date time.Time) (Subscription, error) {
	if date.IsZero() {
		return Subscription{}, httperror.New(400, "date can't be zero")
	}
	city, movie, cinema, err := getMovieCityAndCinema(cityID, movieID, cinemaID)
	if err != nil {
		return Subscription{}, err
	}
	return Subscription{
		ID:        utils.RandomString(5),
		CityID:    city.ID,
		MovieID:   movie.ID,
		CinemaID:  cinema.ID,
		CreatedAt: db.UnixTime{Time: time.Now()},
	}, nil
}

func getMovieCityAndCinema(cityID, movieID, cinemaID string) (cities.City, movies.Movie, cinemas.Cinema, error) {
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
	cityOutput := cities.FindByID(ctx, cityID)
	movieOutput := movies.FindByID(ctx, movieID)
	cinemaOutput := cinemas.FindByID(ctx, cinemaID)

	var city cities.City
	var movie movies.Movie
	var cinema cinemas.Cinema
	fmt.Println("started to wait")
	for cityOutput != nil || movieOutput != nil || cinemaOutput != nil {
		fmt.Printf("for loop, %v, %v, %v\n", cityOutput, movieOutput, cinemaOutput)
		//noinspection GoNilness
		select {
		case result := <-cityOutput:
			fmt.Printf("city result: %+v\n", result)
			if result.Err != nil {
				cancel()
				fmt.Printf("city error: %+v\n", result.Err)
				return cities.City{}, movies.Movie{}, cinemas.Cinema{}, result.Err
			}
			city = result.City
			cityOutput = nil
		case result := <-movieOutput:
			fmt.Printf("movie result: %+v\n", result)
			if result.Err != nil {
				cancel()
				fmt.Printf("movie error: %+v\n", result.Err)
				return cities.City{}, movies.Movie{}, cinemas.Cinema{}, result.Err
			}
			movie = result.Movie
			movieOutput = nil
		case result := <-cinemaOutput:
			fmt.Printf("cinema result: %+v\n", result)
			if result.Err != nil {
				cancel()
				fmt.Printf("cinema error: %+v\n", result.Err)
				return cities.City{}, movies.Movie{}, cinemas.Cinema{}, result.Err
			}
			cinema = result.Cinema
			cinemaOutput = nil
		}
	}
	return city, movie, cinema, nil
}
