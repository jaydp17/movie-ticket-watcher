package cinemas

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
	"log"
	"sync"
)

// Fetch cinemas from all the providers & merge them
func Fetch(city cities.City) []Cinema {
	bmsProvider := bookmyshow.Provider{}
	ptmProvider := paytm.Provider{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var bmsCinemas []providers.Cinema
	go func() {
		var err error
		_, bmsCinemas, err = bmsProvider.FetchMoviesAndCinemas(city.BookmyshowID)
		if err != nil {
			log.Printf("error fetching cinemas from bms: %v", err)
		}
		wg.Done()
	}()

	var ptmCinemas []providers.Cinema
	go func() {
		var err error
		_, ptmCinemas, err = ptmProvider.FetchMoviesAndCinemas(city.PaytmID)
		if err != nil {
			log.Printf("error fetching cinemas from ptm: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	mergedCinemas := Merge(bmsCinemas, ptmCinemas)
	return mergedCinemas
}
