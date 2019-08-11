package bookmyshow

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProvider_FetchAvailableVenueCodes(t *testing.T) {
	t.Run("verify args to BMS API", func(t *testing.T) {
		bmsCityID := "BANG"
		bmsChildEventID := "ET00108257"
		date := db.YYYYMMDDTime{Time: time.Now()}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			queryParams := req.URL.Query()
			assert.Equal(t, queryParams.Get("regionCode"), bmsCityID, "incorrect bmsCityID passed to the BMS API")
			assert.Equal(t, queryParams.Get("eventCode"), bmsChildEventID, "incorrect bmsChildEventID passed to BMS API")
			assert.Equal(t, queryParams.Get("dateCode"), date.Format("20060102"), "incorrect date format passed to BMS API, it has to be YYYYMMDD no space/seperator between all 8 characters")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ShowDetails":[{"Venues":[{"VenueCode":"CFBS"}]}]}`))
		}))
		defer server.Close()

		p := Provider{urlToFetchShowTimings: server.URL}
		<-p.FetchAvailableVenueCodes(bmsCityID, bmsChildEventID, date)
	})

	t.Run("verify parsed venue codes", func(t *testing.T) {
		bmsCityID := "BANG"
		bmsChildEventID := "ET00108257"
		date := db.YYYYMMDDTime{Time: time.Now()}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ShowDetails":[{"Venues":[{"VenueCode":"CFBS"},{"VenueCode":"INMB"},{"VenueCode":"INRZ"},{"VenueCode":"PVEG"}]}]}`))
		}))
		defer server.Close()

		p := Provider{urlToFetchShowTimings: server.URL}
		venueCodesResult := <-p.FetchAvailableVenueCodes(bmsCityID, bmsChildEventID, date)
		assert.Len(t, venueCodesResult.Data, 4)
		assert.Equal(t, venueCodesResult.Data, []string{"CFBS", "INMB", "INRZ", "PVEG"})
	})
}
