package dao

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

type Cinema struct {
	providers.Cinema
	BookmyshowID string `json:"bookmyshowID,omitempty"`
	PaytmID      string `json:"paytmID,omitempty"`
}

type Cinemas []Cinema

var CinemaTableName = db.GetFullTableName("cinemas")

func (c Cinema) HasAllProviderIDs() bool {
	return len(c.BookmyshowID) > 0 && len(c.PaytmID) > 0
}
