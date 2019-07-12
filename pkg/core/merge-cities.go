package core

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/dao"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

// MergeCities merges the cities obtained from the providers
func MergeCities(bmsCities, pytmCities []providers.City) []dao.City {
	citiesMap := make(map[string]dao.City)

	// bookmyshow cities
	for _, city := range bmsCities {
		citiesMap[city.Name] = dao.City{
			ID:           city.ID,
			Name:         city.Name,
			BookmyshowID: city.ID,
			IsTopCity:    city.IsTopCity,
		}
	}

	for _, city := range pytmCities {
		if existingCity, ok := citiesMap[city.Name]; ok {
			existingCity.PaytmID = city.ID
			citiesMap[city.Name] = existingCity
		}
		// we don't add cities that are not in BookMyShow
		// and are in PayTM
	}

	var cities []dao.City
	for _, city := range citiesMap {
		cities = append(cities, city)
	}
	return cities
}
