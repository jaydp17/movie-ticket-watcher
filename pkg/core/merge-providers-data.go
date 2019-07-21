package core

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/dao"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/utils"
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
	maxMovies := utils.MaxInt(len(bmsMovies), len(pytmMovies))
	moviesMap := make(map[string]dao.Movie, maxMovies)

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

func MergeCinemas(bmsCinemas, pytmCinemas []providers.Cinema) dao.Cinemas {
	cinemasMergedByName, bmsRemaining1, ptmRemaining1 := mergeCinemasByName(bmsCinemas, pytmCinemas)
	cinemasMergedByDistance, bmsRemaining2, ptmRemaining2 := mergeCinemasByGeoDistance(bmsRemaining1, ptmRemaining1)

	allMergedCinemas := make([]dao.Cinema, 0, len(cinemasMergedByName)+len(cinemasMergedByDistance))
	allMergedCinemas = append(allMergedCinemas, cinemasMergedByName...)
	allMergedCinemas = append(allMergedCinemas, cinemasMergedByDistance...)

	bmsFinalConverted := make([]dao.Cinema, 0, len(bmsRemaining2))
	for _, bCinema := range bmsRemaining2 {
		c := dao.Cinema{
			Cinema: providers.Cinema{
				ID:        bCinema.NameSlug(),
				Name:      bCinema.Name,
				Provider:  bCinema.Provider,
				CityID:    bCinema.CityID,
				Latitude:  bCinema.Latitude,
				Longitude: bCinema.Longitude,
				Address:   bCinema.Address,
			},
			BookmyshowID: bCinema.ID,
		}
		bmsFinalConverted = append(bmsFinalConverted, c)
	}
	allMergedCinemas = append(allMergedCinemas, bmsFinalConverted...)

	ptmFinalConverted := make([]dao.Cinema, 0, len(ptmRemaining2))
	for _, pCinema := range ptmRemaining2 {
		c := dao.Cinema{
			Cinema: providers.Cinema{
				ID:        pCinema.NameSlug(),
				Name:      pCinema.Name,
				Provider:  pCinema.Provider,
				CityID:    pCinema.CityID,
				Latitude:  pCinema.Latitude,
				Longitude: pCinema.Longitude,
				Address:   pCinema.Address,
			},
			PaytmID: pCinema.ID,
		}
		ptmFinalConverted = append(ptmFinalConverted, c)
	}
	allMergedCinemas = append(allMergedCinemas, ptmFinalConverted...)

	return allMergedCinemas
}
