package paytm

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"strconv"
	"strings"
)

const macOsUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"

func (p Provider) FetchMoviesAndCinemas(ptmCityID string) ([]providers.Movie, []providers.Cinema, error) {
	params := req.Param{
		"groupResult": "true",
		"city":        ptmCityID,
		"channel":     "web",
		"version":     "2",
	}
	headers := req.Header{"User-Agent": macOsUserAgent}
	res, err := req.Get("https://apiproxy-moviesv2.paytm.com/v2/movies/search", params, headers)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch PayTM Cinemas & movies: %v", err)
	}
	var jsonRes paytmSearchResponse
	if err := res.ToJSON(&jsonRes); err != nil {
		return nil, nil, fmt.Errorf("failed to parse response from PayTM: %v", err)
	}

	var movies []providers.Movie
	for _, ptmMovie := range jsonRes.Movies {
		for _, groupChild := range ptmMovie.Grouped {
			groupID := ptmMovie.generateGroupID()
			movies = append(movies, providers.Movie{
				ID:           groupChild.MovieCode,
				GroupID:      groupID,
				Title:        ptmMovie.Title,
				ScreenFormat: groupChild.Properties.getScreenFormat(),
				Language:     ptmMovie.Language,
				ImageURL:     ptmMovie.ImageURL,
			})
		}
	}

	var cinemas []providers.Cinema
	for _, ptmCinema := range jsonRes.Cinemas {
		cinemas = append(cinemas, providers.Cinema{
			ID:        strconv.Itoa(ptmCinema.ID),
			Name:      ptmCinema.Name,
			Provider:  ptmCinema.ProviderChain,
			CityID:    "", // this is supposed to be the common cityID, which we don't have here, but will be filled later in merging phase
			Latitude:  ptmCinema.Latitude,
			Longitude: ptmCinema.Longitude,
			Address:   ptmCinema.Address,
		})
	}
	return movies, cinemas, nil
}

type paytmSearchResponse struct {
	Cinemas map[string]paytmCinema `json:"cinemas"`
	Movies  map[string]paytmMovie  `json:"movies"`
}

type paytmMovie struct {
	ImageURL     string            `json:"image_url"`
	Display      string            `json:"display"`      // eg. "The Lion King"
	Title        string            `json:"title"`        // eg. "The Lion King"
	Language     string            `json:"language"`     // eg. "English"
	OpeningDate  string            `json:"openingDate"`  // eg. "2019-07-19T00:00:00Z"
	FirstSession string            `json:"firstSession"` // eg. "2019-07-20T08:30:00.000Z"
	WebImgURL    string            `json:"web_img_url"`
	Rank         float32           `json:"rank"` // eg. 0.17
	EdgeBanner   string            `json:"edge_banner"`
	Grouped      []paytmMovieGroup `json:"grouped"`
}

func (m paytmMovie) generateGroupID() string {
	return strings.ReplaceAll(strings.ToLower(m.Title), " ", "-")
}

type paytmMovieGroup struct {
	MovieCode  string               `json:"movieCode"` // eg. "O9QJZ5"
	Properties paytmMovieProperties `json:"properties"`
}

type paytmMovieProperty struct {
	Key     string `json:"key"`   // eg. "screenFormat"
	Value   string `json:"value"` // eg. "IMAX 3D"
	Display bool   `json:"display"`
}
type paytmMovieProperties []paytmMovieProperty

func (props paytmMovieProperties) getScreenFormat() string {
	for _, prop := range props {
		if prop.Key == "screenFormat" {
			return prop.Value
		}
	}
	return ""
}

type paytmCinema struct {
	ID               int     `json:"id"`            // eg. 12
	Name             string  `json:"name"`          // eg. "PVR Forum Mall, Koramangala"
	ProviderId       int     `json:"providerId"`    // eg. 1
	ProviderName     string  `json:"providerName"`  // eg. "PVR"
	ProviderType     string  `json:"providerType"`  // eg. "provider"
	ProviderChain    string  `json:"providerChain"` // eg. "PVR"
	Address          string  `json:"address"`       // eg. "The Forum Mall, 21-22, Adugodi Main Road, Chikku Lakshmaiah Layout, Koramangala, Bengaluru, Karnataka 560095, India"
	SubCity          string  `json:"subCity"`       // eg. "Bengaluru"
	Latitude         float64 `json:"latitude"`      // eg. 12.934661
	Longitude        float64 `json:"longitude"`     // eg. 77.611314
	CoverImageURL    string  `json:"coverImageUrl"`
	AppCoverImageURL string  `json:"appCoverImageUrl"`
}
