package cinemas

import "github.com/jaydp17/movie-ticket-watcher/pkg/config"

type Cinema struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Provider     string  `json:"provider"`
	CityID       string  `json:"cityID"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Address      string  `json:"address"`
	BookmyshowID string  `json:"bookmyshowID,omitempty"`
	PaytmID      string  `json:"paytmID,omitempty"`
}

var TableName = config.FullTableName("cinemas")

func (c Cinema) HasAllProviderIDs() bool {
	return len(c.BookmyshowID) > 0 && len(c.PaytmID) > 0
}
