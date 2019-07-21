package core

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/dao"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

// MergeCities merges the cities obtained from the providers
func MergeCities(bmsCities, pytmCities []providers.City) dao.Cities {
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
	}
	// we don't add cities that are not in both providers
	var cities dao.Cities
	for _, city := range citiesMap {
		if city.HasAllProviderIDs() {
			cities = append(cities, city)
		}
	}
	return cities
}

// MergeMovies merges the movies obtained from the providers
func MergeMovies(bmsMovies, pytmMovies []providers.Movie) dao.Movies {
	moviesMap := make(map[string]dao.Movie)

	// bookmyshow movies
	for _, movie := range bmsMovies {
		slug := movie.Slug()
		movieWithBmsID := dao.Movie{
			Movie: providers.Movie{
				ID:           slug,
				GroupID:      movie.GroupSlug(),
				Title:        movie.Title,
				ScreenFormat: movie.ScreenFormat,
				Language:     movie.Language,
				ImageURL:     movie.ImageURL,
			},
			BookmyshowID: movie.ID,
		}
		moviesMap[slug] = movieWithBmsID
	}

	// paytm movies
	for _, movie := range pytmMovies {
		slug := movie.Slug()
		if existingMovie, ok := moviesMap[slug]; ok {
			existingMovie.PaytmID = movie.ID
			// the execution comes here only if the title, screenFormat & the language are the same in both the providers
			moviesMap[slug] = existingMovie
		} else {
			movieWithPtmID := dao.Movie{
				Movie: providers.Movie{
					ID:           slug,
					GroupID:      movie.GroupSlug(),
					Title:        movie.Title,
					ScreenFormat: movie.ScreenFormat,
					Language:     movie.Language,
					ImageURL:     movie.ImageURL,
				},
				PaytmID: movie.ID,
			}
			moviesMap[slug] = movieWithPtmID
		}
	}

	mergedMovies := make([]dao.Movie, 0, len(moviesMap))
	for _, m := range moviesMap {
		mergedMovies = append(mergedMovies, m)
	}
	return mergedMovies
}
