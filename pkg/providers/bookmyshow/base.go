package bookmyshow

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

type Provider struct {
	providers.Provider
	urlToFetchCities           string
	urlToFetchMoviesAndCinemas string
	urlToFetchShowTimings      string
}

type YesNo string

const (
	okHTTPUserAgent = "okhttp/3.11.0"
	token           = "67x1xa33b4x422b361ba"
)

func New() Provider {
	return Provider{
		urlToFetchCities:           "https://data-in.bookmyshow.com",
		urlToFetchMoviesAndCinemas: "https://in.bookmyshow.com/serv/getData",
		urlToFetchShowTimings:      "https://in.bookmyshow.com/api/v2/mobile/showtimes/byevent",
	}
}
