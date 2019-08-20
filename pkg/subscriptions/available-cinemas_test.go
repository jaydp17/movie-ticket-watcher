package subscriptions

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockProviderForAreTicketsAvailable struct {
	result providers.VenueCodesResult
}

func (p mockProviderForAreTicketsAvailable) FetchAvailableVenueCodes(providerCityID, providerMovieID string, date db.YYYYMMDDTime) <-chan providers.VenueCodesResult {
	ch := make(chan providers.VenueCodesResult)
	go func() {
		ch <- p.result
		close(ch)
	}()
	return ch
}

func TestAreTicketsAvailable(t *testing.T) {
	type VenueCodesResult = providers.VenueCodesResult
	city := cities.City{BookmyshowID: "BANG", PaytmID: "bengaluru"}
	movie := movies.Movie{BookmyshowID: "MX012345", PaytmID: "QCB12345"}
	cinema := cinemas.Cinema{BookmyshowID: "PVRMX", PaytmID: "123"}

	tests := []struct {
		name      string
		bmsResult VenueCodesResult
		ptmResult VenueCodesResult
		date      db.YYYYMMDDTime
		result    bool
		err       error
	}{
		{name: "available in both providers", bmsResult: VenueCodesResult{Data: []string{cinema.BookmyshowID}}, ptmResult: VenueCodesResult{Data: []string{cinema.PaytmID}}, result: true},
		{name: "available in just bms", bmsResult: VenueCodesResult{Data: []string{cinema.BookmyshowID}}, ptmResult: VenueCodesResult{Data: []string{}}, result: true},
		{name: "available in just ptm", bmsResult: VenueCodesResult{Data: []string{}}, ptmResult: VenueCodesResult{Data: []string{cinema.PaytmID}}, result: true},
		{name: "available in no provider", bmsResult: VenueCodesResult{Data: []string{}}, ptmResult: VenueCodesResult{Data: []string{}}, result: false},
		{name: "error in bms request", bmsResult: VenueCodesResult{Err: fmt.Errorf("request blocked")}, ptmResult: VenueCodesResult{Data: []string{cinema.PaytmID}}, result: true},
		{name: "error in ptm request", bmsResult: VenueCodesResult{Data: []string{cinema.BookmyshowID}}, ptmResult: VenueCodesResult{Err: fmt.Errorf("request blocked")}, result: true},
		{name: "error in both providers request", bmsResult: VenueCodesResult{Err: fmt.Errorf("request blocked")}, ptmResult: VenueCodesResult{Err: fmt.Errorf("request blocked")}, result: false},
		{name: "error on past date", result: false, err: DateInPastError{db.YYYYMMDDFromTime(time.Now().AddDate(0, 0, -2))}, date: db.YYYYMMDDFromTime(time.Now().AddDate(0, 0, -2))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bmsProvider := mockProviderForAreTicketsAvailable{result: tt.bmsResult}
			ptmProvider := mockProviderForAreTicketsAvailable{result: tt.ptmResult}

			date := tt.date
			if date.IsZero() {
				date = db.YYYYMMDDFromTime(time.Now())
			}
			got, err := AreTicketsAvailable(bmsProvider, ptmProvider, city, movie, cinema, date)
			assert.Equal(t, tt.result, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
