package bookmyshow

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"github.com/jaydp17/movie-ticket-watcher/pkg/utils"
	"strings"
)

func (p Provider) FetchMoviesAndCinemas(bmsCityID string) ([]providers.Movie, []providers.Cinema, error) {
	params := req.Param{
		"cmd":                "QUICKBOOK",
		"type":               "MT",
		"getRecommendedData": "1",
	}
	headers := req.Header{
		"User-Agent": okHTTPUserAgent,
		"Cookie":     fmt.Sprintf("Rgn=|Code=%s", bmsCityID),
	}
	res, err := req.Get(p.urlToFetchMoviesAndCinemas, params, headers)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch BMS Cinemas & movies: %v", err)
	}
	var jsonRes bmsQuickBookResponse
	if err := res.ToJSON(&jsonRes); err != nil {
		return nil, nil, fmt.Errorf("failed to parse response from BMS: %v", err)
	}

	var movies []providers.Movie
	for _, bmsMovie := range jsonRes.MoviesData.BookMyShow.ArrEvents {
		for _, movieChild := range bmsMovie.ChildEvents {
			groupID := bmsMovie.generateGroupID()
			movies = append(movies, providers.Movie{
				ID:           movieChild.EventCode,
				GroupID:      groupID,
				Title:        bmsMovie.EventTitle,
				Language:     movieChild.EventLanguage,
				ScreenFormat: movieChild.screenFormat(),
				ImageURL:     movieChild.getImageURL(),
			})
		}
	}

	var cinemas []providers.Cinema
	for _, bmsCinema := range jsonRes.Cinemas.BookMyShow.AiVN {
		lat := utils.ToFloat(bmsCinema.VenueLatitude)
		lng := utils.ToFloat(bmsCinema.VenueLongitude)
		cinemas = append(cinemas, providers.Cinema{
			ID:        bmsCinema.VenueCode,
			Name:      bmsCinema.VenueName,
			Provider:  bmsCinema.CompanyCode,
			CityID:    "", // this is supposed to be the common cityID, which we don't have here, but will be filled later in merging phase
			Latitude:  lat,
			Longitude: lng,
			Address:   bmsCinema.VenueAddress,
		})
	}

	return movies, cinemas, nil
}

type bmsQuickBookResponse struct {
	MoviesData struct {
		BookMyShow struct {
			ArrEvents []bmsQuickBookMovie `json:"arrEvents"`
		}
	} `json:"moviesData"`
	Cinemas struct {
		BookMyShow struct {
			AiVN []bmsQuickBookCinema `json:"aiVN"`
		}
	} `json:"cinemas"`
}

type bmsQuickBookMovie struct {
	EventCode          string                  // eg. 'ET00056555'
	EventGroup         string                  // eg. 'EG00037499'
	EventTitle         string                  // eg. 'Captain Marvel'
	EventGrpDuration   string                  // eg. '124'
	EventGrpSequence   string                  // eg. '50'
	EventGrpGenre      string                  // eg. '|ACTION|ADVN|FANTASY|'
	EventGrpCensor     string                  // eg. 'UA'
	EventGrpIsWebView  string                  // eg. 'false'
	IsMovieClubEnabled YesNo                   // eg. 'Y' | 'N'
	EventURLTitle      string                  // eg. 'captain-marvel'
	IsPremiere         YesNo                   // eg. 'Y' | 'N'
	AvgRating          int `json:"avgRating"`  // eg. 83
	TotalVotes         int `json:"totalVotes"` // eg. 48011
	ChildEvents        []bmsQuickBookMovieChild
}

func (m bmsQuickBookMovie) generateGroupID() string {
	return strings.ReplaceAll(strings.ToLower(m.EventTitle), " ", "-")
}

type bmsQuickBookMovieChild struct {
	EventCode           string // eg. 'ET00100559'
	EventType           string // eg. 'MT'
	EventLanguage       string // eg. 'English'
	EventName           string // eg. 'Avengers: Endgame (3D)'
	EventGenre          string // eg. 'Action|Adventure|Fantasy'
	EventDimension      string // eg. '3D' or 'IMAX 3D'
	EventImageCode      string // eg. 'the-lion-king-et00105989-28-06-2019-02-55-52' which needs to be converted to https://in.bmscdn.com/iedb/movies/images/mobile/thumbnail/xlarge/the-lion-king-et00105989-28-06-2019-02-55-52.jpg
	EventIsAtmosEnabled YesNo
}

func (movieChild bmsQuickBookMovieChild) screenFormat() string {
	cleanFormat := utils.KeepJustAlphaNumeric(strings.ToLower(movieChild.EventDimension))
	if cleanFormat == "4dx3d" {
		return "4dx"
	}
	return cleanFormat
}
func (movieChild bmsQuickBookMovieChild) getImageURL() string {
	if len(movieChild.EventImageCode) > 0 {
		return fmt.Sprintf("https://in.bmscdn.com/iedb/movies/images/mobile/thumbnail/xlarge/%s.jpg", movieChild.EventImageCode)
	}
	return ""
}

type bmsQuickBookCinema struct {
	VenueCode          string // eg. 'ENNR'
	CompanyCode        string // eg. 'MOAM'
	VenueName          string // eg. '7D Mastii: Element Mall, Nagwara'
	IsATMOSEnabled     YesNo
	VenueLatitude      string // eg. '13.0452'
	VenueLongitude     string // eg. '77.6266'
	VenueAddress       string // eg. '3rd floor, Nagavara Village, 100 wide Thanisandra main road, Bengaluru, Karnataka 560077, India'
	VenueSubRegionCode string // eg. 'BANG'
	VenueSubRegionName string // eg. 'Bengaluru'
	CinemaUnpaidFlag   YesNo
	VenueLegends       string // eg. ';CAR;FOD;';
	CinemaIsNew        YesNo
	CinemaCodFlag      YesNo
	CinemaCopFlag      YesNo
}
