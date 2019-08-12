package paytm

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestProvider_FetchAvailableVenueCodes(t *testing.T) {
	t.Run("verify args to PTM API", func(t *testing.T) {
		ptmCityID := "bengaluru"
		ptmMovieID := "O9QTPF"
		date := db.YYYYMMDDTime{Time: time.Now()}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			queryParams := req.URL.Query()
			assert.Equal(t, queryParams.Get("city"), ptmCityID, "incorrect city passed to the PTM API")
			assert.Equal(t, queryParams.Get("moviecode"), ptmMovieID, "incorrect moviecode passed to PTM API")
			assert.Equal(t, queryParams.Get("fromdate"), date.ToYYYYMMDD(), "incorrect date format passed to PTM API, it has to be YYYY-MM-DD")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(strconv.Quote(`{"movies":{"O9QTPF":{"sessions":[{"movieCode":"O9QTPF","cinemaId":305},{"movieCode":"O9QTPF","cinemaId":305},{"movieCode":"O9QTPF","cinemaId":305}]}}}"`)))
		}))
		defer server.Close()

		p := Provider{urlToFetchShowTimings: server.URL}
		<-p.FetchAvailableVenueCodes(ptmCityID, ptmMovieID, date)
	})

	t.Run("verify parsed venue codes", func(t *testing.T) {
		ptmCityID := "bengaluru"
		ptmMovieID := "O9QTPF"
		date := db.YYYYMMDDTime{Time: time.Now()}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(strconv.Quote(`{"movies":{"O9QTPF":{"sessions":[{"movieCode":"O9QTPF","cinemaId":305},{"movieCode":"O9QTPF","cinemaId":305},{"movieCode":"O9QTPF","cinemaId":3674}]}}}`)))
		}))
		defer server.Close()

		p := Provider{urlToFetchShowTimings: server.URL}
		venueCodesResult := <-p.FetchAvailableVenueCodes(ptmCityID, ptmMovieID, date)
		assert.Equal(t, venueCodesResult.Data, []string{"305", "3674"})
	})

	t.Run("capture error", func(t *testing.T) {
		ptmCityID := "bengalur" // <-- there's an intentional typo here
		ptmMovieID := "O9QTPF"
		date := db.YYYYMMDDTime{Time: time.Now()}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(strconv.Quote(`{"error":"Invalid City field.(search)","status":{"result":"failure","message":{"title":"City not found","message":"Invalid City field.(search)"}},"code":404,"mCode":"CITY_NF"}`)))
		}))
		defer server.Close()

		p := Provider{urlToFetchShowTimings: server.URL}
		venueCodesResult := <-p.FetchAvailableVenueCodes(ptmCityID, ptmMovieID, date)
		assert.Equal(t, venueCodesResult.Err, fmt.Errorf("error from PayTM: Invalid City field.(search)"))
	})
}
