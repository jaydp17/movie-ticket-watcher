package movies

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/moviecitylink"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
	"log"
	"sync"
)

// Fetch movies from all the providers & merge them
func Fetch(city cities.City) []Movie {
	movies := fetchAndMerge(city)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Write movies
	go func() {
		if writeErr := Write(movies); writeErr != nil {
			log.Printf("error writing movies to db: %+v", writeErr)
		}
		wg.Done()
	}()

	// Write movie-city-link
	go func() {
		movieIDs := make([]string, 0, len(movies))
		for _, m := range movies {
			movieIDs = append(movieIDs, m.ID)
		}
		link := moviecitylink.MovieCityLink{
			CityID:   city.ID,
			MovieIDs: movieIDs,
		}
		if writeErr := moviecitylink.WriteOne(link); writeErr != nil {
			log.Printf("error writing movieCityLink to db: %+v", writeErr)
		}
		wg.Done()
	}()

	wg.Wait()
	return movies
}

func fetchAndMerge(city cities.City) []Movie {
	bmsProvider := bookmyshow.Provider{}
	ptmProvider := paytm.Provider{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var bmsMovies []providers.Movie
	go func() {
		var err error
		bmsMovies, _, err = bmsProvider.FetchMoviesAndCinemas(city.BookmyshowID)
		if err != nil {
			log.Printf("error fetching movies from bms: %v", err)
		}
		wg.Done()
	}()

	var ptmCinemas []providers.Movie
	go func() {
		var err error
		ptmCinemas, _, err = ptmProvider.FetchMoviesAndCinemas(city.PaytmID)
		if err != nil {
			log.Printf("error fetching movies from ptm: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	mergedCinemas := Merge(bmsMovies, ptmCinemas)
	return mergedCinemas
}
