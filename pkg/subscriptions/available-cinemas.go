package subscriptions

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
)

func AreTicketsAvailable(city cities.City, movie movies.Movie, cinema cinemas.Cinema, date db.YYYYMMDDTime) (bool, error) {
	bms := bookmyshow.Provider{}
	var bmsVenuesResult <-chan providers.VenueCodesResult
	if len(cinema.BookmyshowID) > 0 {
		bmsVenuesResult = bms.FetchAvailableVenueCodes(city.BookmyshowID, movie.BookmyshowID, date)
	}

	ptm := paytm.Provider{}
	var ptmVenuesResult <-chan providers.VenueCodesResult
	if len(cinema.PaytmID) > 0 {
		ptmVenuesResult = ptm.FetchAvailableVenueCodes(city.PaytmID, movie.PaytmID, date)
	}

	bmsVenueCodes := make([]string, 0)
	ptmVenueCodes := make([]string, 0)
	for bmsVenuesResult != nil || ptmVenuesResult != nil {
		//noinspection GoNilness
		select {
		case result := <-bmsVenuesResult:
			if result.Err == nil && len(result.Data) > 0 {
				bmsVenueCodes = result.Data
			}
			if result.Err != nil {
				fmt.Println(result.Err)
			}
			bmsVenuesResult = nil
		case result := <-ptmVenuesResult:
			if result.Err == nil && len(result.Data) > 0 {
				ptmVenueCodes = result.Data
			}
			if result.Err != nil {
				fmt.Println(result.Err)
			}
			ptmVenuesResult = nil
		}
	}

	for _, bmsVenueCode := range bmsVenueCodes {
		if cinema.BookmyshowID == bmsVenueCode {
			return true, nil
		}
	}
	for _, ptmVenueCode := range ptmVenueCodes {
		if cinema.PaytmID == ptmVenueCode {
			return true, nil
		}
	}
	return false, nil
}
