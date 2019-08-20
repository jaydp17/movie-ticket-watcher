package subscriptions

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"time"
)

func AreTicketsAvailable(bms providers.AvailableVenueCodesFetcher, ptm providers.AvailableVenueCodesFetcher, city cities.City, movie movies.Movie, cinema cinemas.Cinema, date db.YYYYMMDDTime) (bool, error) {
	today := getBeginningOfDay(time.Now())
	requestedDate := getBeginningOfDay(date.Time)
	if requestedDate.Before(today) {
		return false, DateInPastError{date}
	}

	var bmsVenuesResult <-chan providers.VenueCodesResult
	if len(cinema.BookmyshowID) > 0 {
		bmsVenuesResult = bms.FetchAvailableVenueCodes(city.BookmyshowID, movie.BookmyshowID, date)
	}

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

func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

type DateInPastError struct {
	date db.YYYYMMDDTime
}

func (d DateInPastError) Error() string {
	return fmt.Sprintf("can't check tickets in the past (%s)", d.date.ToYYYYMMDD())
}
