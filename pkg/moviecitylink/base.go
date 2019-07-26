package moviecitylink

import "github.com/jaydp17/movie-ticket-watcher/pkg/config"

type MovieCityLink struct {
	CityID   string   `json:"cityID"`
	MovieIDs []string `json:"movieIDs"`
}

var TableName = config.FullTableName("movie-city-link")
