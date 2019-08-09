package paytm

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"strconv"
)

func (p Provider) GetAvailableVenueCodes(ptmCityID, ptmMovieID string, date db.YYYYMMDDTime) ([]string, error) {
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

	availableVenueCodes := make([]string, 0, len(jsonRes.Cinemas))
	for _, cinema := range jsonRes.Cinemas {
		availableVenueCodes = append(availableVenueCodes, strconv.Itoa(cinema.ID))
	}
	return availableVenueCodes, nil
}

type ptmShowTimeResponse struct {
	Cinemas map[string]paytmCinema `json:"cinemas"`
}
