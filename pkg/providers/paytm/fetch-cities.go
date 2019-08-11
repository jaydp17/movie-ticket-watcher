package paytm

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

const androidUserAgent = "Dalvik/2.1.0 (Linux; U; Android 5.0; Google Build/LRX21M)"

func (p Provider) FetchCities() ([]providers.City, error) {
	headers := req.Header{
		"User-Agent": androidUserAgent,
	}
	res, err := req.Get(p.urlToFetchCities, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PayTM cities: %v", err)
	}
	var jsonRes []payTmCity
	if err := res.ToJSON(&jsonRes); err != nil {
		return nil, fmt.Errorf("failed to parse response from PayTM: %v", err)
	}

	var cities []providers.City
	for _, city := range jsonRes {
		cities = append(cities, providers.City{
			ID:        city.Value,
			Name:      city.Label,
			IsTopCity: city.IsTopCity,
		})
	}
	return cities, nil
}

type payTmCity struct {
	Label     string
	Value     string
	IsTopCity bool
}
