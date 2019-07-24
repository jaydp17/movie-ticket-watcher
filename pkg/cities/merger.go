package cities

import "github.com/jaydp17/movie-ticket-watcher/pkg/providers"

// Merge the cities obtained from the providers
func Merge(bmsCities, pytmCities []providers.City) []City {
	citiesMap := make(map[string]City)

	// bookmyshow cities
	for _, city := range bmsCities {
		citiesMap[city.Name] = City{
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
	}
	// we don't add cities that are not in both providers
	var cities []City
	for _, city := range citiesMap {
		if city.HasAllProviderIDs() {
			cities = append(cities, city)
		}
	}
	return cities
}
