package paytm

import "github.com/jaydp17/movie-ticket-watcher/pkg/providers"

type Provider struct {
	providers.Provider
	urlToFetchCities           string
	urlToFetchMoviesAndCinemas string
	urlToFetchShowTimings      string
}

func New() Provider {
	return Provider{
		urlToFetchCities:           "https://tickets.paytm.com/v1/movies/cities",
		urlToFetchMoviesAndCinemas: "https://apiproxy-moviesv2.paytm.com/v2/movies/search",
		urlToFetchShowTimings:      "https://paytm.com/v1/api/movies/search",
	}
}
