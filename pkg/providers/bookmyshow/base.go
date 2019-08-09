package bookmyshow

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

type Provider struct {
	providers.Provider
}

type YesNo string

const (
	okHTTPUserAgent = "okhttp/3.11.0"
	token           = "67x1xa33b4x422b361ba"
)
