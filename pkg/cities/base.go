package cities

import "github.com/jaydp17/movie-ticket-watcher/pkg/config"

type City struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	BookmyshowID string `json:"bookmyshowID,omitempty"`
	PaytmID      string `json:"paytmID,omitempty"`
	IsTopCity    bool   `json:"isTopCity"`
}

var TableName = config.FullTableName("cities")

// HasAllProviderIDs checks if a cities object has all the provider IDs
func (c City) HasAllProviderIDs() bool {
	return len(c.BookmyshowID) > 0 && len(c.PaytmID) > 0
}
