package movies

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
	"log"
	"sync"
)

// Fetch movies from all the providers & merge them
func Fetch(city cities.City) []Movie {
	movies := fetchAndMerge(city)

	writeErr := Write(movies)
	if writeErr != nil {
		log.Printf("error writing movies to db: %+v", writeErr)
	}
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
