package paytm

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"strconv"
)

func (p Provider) FetchAvailableVenueCodes(ptmCityID, ptmMovieID string, date db.YYYYMMDDTime) ([]string, error) {
	params := req.Param{
		"groupResult": "true",
		"city":        ptmCityID,
		"moviecode":   ptmMovieID,
		"fromdate":    date.ToYYYYMMDD(),
		"channel":     "web",
		"version":     "2",
	}
	headers := req.Header{"User-Agent": macOsUserAgent}
	res, err := req.Get("https://paytm.com/v1/api/movies/search", params, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch showtimes from PayTM: %v", err)
	}

	strResponse, err := res.ToString()
	if err != nil {
		return nil, fmt.Errorf("failed to get string response from PayTM: %v", err)
	}

	jsonStr, err := strconv.Unquote(strResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get unquoted string: %v", err)
	}

	var jsonRes ptmShowTimeResponse
	if err := json.Unmarshal([]byte(jsonStr), &jsonRes); err != nil {
		return nil, err
	}
	if len(jsonRes.Error) != 0 {
		return nil, fmt.Errorf("error from PayTM: %v", jsonRes.Error)
	}

	ptmMovie, ok := jsonRes.Movies[ptmMovieID]
	if !ok {
		return nil, fmt.Errorf("invalid response format from PayTM")
	}

	availableVenueCodes := make([]string, 0)
	for _, session := range ptmMovie.Sessions {
		if session.MovieCode == ptmMovieID {
			availableVenueCodes = append(availableVenueCodes, strconv.Itoa(session.CinemaID))
		}
	}
	return availableVenueCodes, nil
}

type ptmShowTimeResponse struct {
	Error   string                            `json:"error"`
	Movies  map[string]paytmMovieWithSessions `json:"movies"`
	Cinemas map[string]paytmCinema            `json:"cinemas"`
}

type paytmMovieWithSessions struct {
	ImageURL     string               `json:"image_url"`
	Display      string               `json:"display"`      // eg. "The Lion King"
	Title        string               `json:"title"`        // eg. "The Lion King"
	Language     string               `json:"language"`     // eg. "English"
	OpeningDate  string               `json:"openingDate"`  // eg. "2019-07-19T00:00:00Z"
	FirstSession string               `json:"firstSession"` // eg. "2019-07-20T08:30:00.000Z"
	WebImgURL    string               `json:"web_img_url"`
	Rank         float32              `json:"rank"` // eg. 0.17
	EdgeBanner   string               `json:"edge_banner"`
	Sessions     []paytmMovieSessions `json:"sessions"`
}

type paytmMovieSessions struct {
	MovieCode    string               `json:"movieCode"`
	CinemaID     int                  `json:"cinemaId"`
	CinemaName   string               `json:"cinemaName"`   // eg. "PVR Forum Mall, Koramangala"
	ProviderId   int                  `json:"providerId"`   // eg. 1
	ProviderName string               `json:"providerName"` // eg. "PVR"
	Address      string               `json:"address"`      // eg. "The Forum Mall, 21-22, Adugodi Main Road, Chikku Lakshmaiah Layout, Koramangala, Bengaluru, Karnataka 560095, India"
	Latitude     float64              `json:"latitude"`     // eg. 12.934661
	Longitude    float64              `json:"longitude"`    // eg. 77.611314
	Properties   paytmMovieProperties `json:"properties"`
}
