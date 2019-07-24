package movies

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/utils"
)

// Merge the movies obtained from the providers
func Merge(bmsMovies, pytmMovies []providers.Movie) []Movie {
	maxMovies := utils.MaxInt(len(bmsMovies), len(pytmMovies))
	moviesMap := make(map[string]Movie, maxMovies)

	// bookmyshow movies
	for _, movie := range bmsMovies {
		slug := movie.Slug()
		movieWithBmsID := Movie{
			ID:           slug,
			GroupID:      movie.GroupSlug(),
			Title:        movie.Title,
			ScreenFormat: movie.ScreenFormat,
			Language:     movie.Language,
			ImageURL:     movie.ImageURL,
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
			movieWithPtmID := Movie{
				ID:           slug,
				GroupID:      movie.GroupSlug(),
				Title:        movie.Title,
				ScreenFormat: movie.ScreenFormat,
				Language:     movie.Language,
				ImageURL:     movie.ImageURL,
				PaytmID:      movie.ID,
			}
			moviesMap[slug] = movieWithPtmID
		}
	}

	mergedMovies := make([]Movie, 0, len(moviesMap))
	for _, m := range moviesMap {
		mergedMovies = append(mergedMovies, m)
	}
	return mergedMovies
}
