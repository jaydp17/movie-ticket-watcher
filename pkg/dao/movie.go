package dao

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

type Movie struct {
	providers.Movie
	BookmyshowID string `json:"bookmyshowID,omitempty"`
	PaytmID      string `json:"paytmID,omitempty"`
}

type Movies []Movie

var MovieTableName = db.GetFullTableName("movies")
