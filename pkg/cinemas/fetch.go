package cinemas

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
	"log"
	"sync"
)

// Fetch cinemas from all the providers & merge them
func Fetch(dbClient dynamodbiface.ClientAPI, city cities.City) []Cinema {
	cinemas := fetchAndMerge(city)

	writeErr := Write(dbClient, cinemas)
	if writeErr != nil {
		log.Printf("error writing cinemas to db: %+v", writeErr)
	}
	return cinemas
}

func fetchAndMerge(city cities.City) []Cinema {
	bmsProvider := bookmyshow.New()
	ptmProvider := paytm.New()

	wg := sync.WaitGroup{}
	wg.Add(2)

	var bmsCinemas []providers.Cinema
	go func() {
		var err error
		_, bmsCinemas, err = bmsProvider.FetchMoviesAndCinemas(city.BookmyshowID)
		if err != nil {
			log.Printf("error fetching cinemas from bms: %v", err)
		}
		for i := range bmsCinemas {
			bmsCinemas[i].CityID = city.ID
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
		for i := range ptmCinemas {
			ptmCinemas[i].CityID = city.ID
		}
		wg.Done()
	}()

	wg.Wait()

	mergedCinemas := Merge(bmsCinemas, ptmCinemas)
	return mergedCinemas
}
