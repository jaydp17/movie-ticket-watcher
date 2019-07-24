package movies

import "github.com/jaydp17/movie-ticket-watcher/pkg/config"

type Movie struct {
	ID           string `json:"id"`
	GroupID      string `json:"groupID"`
	Title        string `json:"title"` // make sure there's no additional information in the title like (3D) or the Language
	ScreenFormat string `json:"screenFormat"`
	Language     string `json:"language"`
	ImageURL     string `json:"imageURL"`
	BookmyshowID string `json:"bookmyshowID"`
	PaytmID      string `json:"paytmID"`
}

var TableName = config.FullTableName("movies")
