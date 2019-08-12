package bookmyshow

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"net/http"
)

// FetchAvailableVenueCodes fetches available cinemas where the given movie is screening on that date
func (p Provider) FetchAvailableVenueCodes(bmsCityID, bmsChildEventID string, date db.YYYYMMDDTime) <-chan providers.VenueCodesResult {
	resultCh := make(chan providers.VenueCodesResult)

	if len(bmsCityID) == 0 || len(bmsChildEventID) == 0 {
		resultCh <- providers.VenueCodesResult{Err: fmt.Errorf("bmsCityID or bmsChildEventID is empty")}
		close(resultCh)
		return resultCh
	}

	go func() {
		defer close(resultCh)
		params := req.Param{
			"regionCode": bmsCityID,
			"eventCode":  bmsChildEventID,
			"dateCode":   date.Format("20060102"),
			"token":      token,
			"bmsId":      "1.82650383.1552055894719", // TODO: figure out what is this?
		}
		headers := req.Header{
			"User-Agent": okHTTPUserAgent,
		}
		res, err := req.Get(p.urlToFetchShowTimings, params, headers)
		if err != nil {
			resultCh <- providers.VenueCodesResult{Err: fmt.Errorf("failed to fetch showtimes from BMS: %v", err)}
			return
		}

		if res.Response().StatusCode != http.StatusOK {
			respBody, err := res.ToString()
			if err != nil {
				resultCh <- providers.VenueCodesResult{Err: fmt.Errorf("failed to fetch showtimes from BMS: %v", err)}
				return
			}
			resultCh <- providers.VenueCodesResult{Err: fmt.Errorf("failed to fetch showtimes from BMS (with Status %s) and body: %v", res.Response().Status, respBody)}
			return
		}

		var jsonRes bmsResponse
		if err := res.ToJSON(&jsonRes); err != nil {
			resultCh <- providers.VenueCodesResult{Err: fmt.Errorf("failed to fetch showtimes from BMS: %v", err)}
			return
		}

		if len(jsonRes.ShowDetails) == 0 {
			// when there's no showDetails
			resultCh <- providers.VenueCodesResult{Data: []string{}}
			return
		}
		showDetail := jsonRes.ShowDetails[0]
		availableVenueCodes := showDetail.getVenueCodes()
		resultCh <- providers.VenueCodesResult{Data: availableVenueCodes}
	}()
	return resultCh
}

type bmsResponse struct {
	ShowDetails []BmsShowDetails
}

// BmsShowDetails holds the event & venues information
type BmsShowDetails struct {
	Date   string
	Event  bmsShowTimeEvent
	Venues []bmsShowTimeVenue
}

func (showDetail BmsShowDetails) getVenueCodes() []string {
	var venueCodes []string
	for _, venue := range showDetail.Venues {
		venueCodes = append(venueCodes, venue.VenueCode)
	}
	return venueCodes
}

type bmsShowTimeEvent struct {
	EventTitle string // eg. 'Avengers: Endgame'
	EventGroup string // eg. 'EG00068832'
}
type bmsShowTimeVenue struct {
	VenueCode           string // eg. 'INMB'
	VenueName           string // eg. 'INOX: Mantri Square, Malleshwaram'
	VenueAdd            string // eg. 'Mantri Square Mall, Sampige Road, Malleswaram, Bengaluru, Karnataka 560052, India'
	VenueApp            string // eg. 'SB'
	SubRegSeq           string // eg. '1'
	CouponIsAllowed     string // eg. 'Y' or 'N'
	AllowSales          string // eg. 'Y' or 'N'
	Lng                 string // eg. '77.5703'
	ShowSeatNo          string // eg. 'Y' or 'N'
	SessCount           string // eg. '124'
	SubRegCode          string // eg. 'BANG'
	SubRegName          string // eg. 'BANG'
	TicketCancellation  string // eg. 'Y' or 'N'
	UnpaidReleaseCutOff string // eg. '1 hr'
	CinemaCodFlag       string // eg. 'Y' or 'N'
	IsFullLayout        string // eg. 'Y' or 'N'
	ETicket             string // eg. 'Y' or 'N'
	MTicket             string // eg. 'Y' or 'N'
	BestSeatsAvail      string // eg. 'Y' or 'N' // 🤔 ?
	CoupleSeats         string // eg. 'Y' or 'N'
	CompCode            string // eg. 'INOX'
	ShowTimes           []bmsShowTime
}
type bmsShowTimeChildEvent struct {
	EventCode           string // eg. 'ET00100559'
	EventTitle          string // eg. 'Avengers: Endgame (3D)'
	EventType           string // eg. 'MT'
	EventLang           string // eg. 'English'
	EventName           string // eg. 'Avengers: Endgame (3D) - English'
	EventGenre          string // eg. 'Action|Adventure|Fantasy'
	EventDimension      string // eg. '3D' or 'IMAX 3D'
	EventIsAtmosEnabled YesNo
}

type bmsShowTime struct {
	ShowDateTime           string   // eg. '201904261045'
	MinPrice               string   // eg. '397.00';
	EventCode              string   // eg. 'ET00100668';
	BestAvailableSeats     int      // eg. 0;
	Availability           string   // eg. 'A';
	ShowTime               string   // eg. '10:45 AM';
	ShowDateCode           string   // eg. '20190426';
	SessionUnpaidFlag      string   // eg. 'Y' or 'N'
	CoupleSeats            string   // eg. 'Y' or 'N'
	SessionUnpaidQuota     string   // eg. '0';
	IsAtmosEnabled         string   // eg. 'Y' or 'N'
	MaxPrice               string   // eg. '397.00';
	ApplicablePriceFilters []string // eg. ['pf5'];
	ShowTimeCode           string   // eg. '1045';
	Categories             []bmsShowTimeCategory
}

type bmsShowTimeCategory struct {
	PercentAvail          string // eg. '1';
	PriceCode             string // eg. 'CL';
	AdditionalData        string // eg. '0';
	CurPrice              string // eg. '397.00';
	AreaCatCode           string // eg. 'CL';
	MaxSeats              string // eg. '246';
	BestAvailableSeats    string // eg. '0';
	SeatLayout            string // eg. 'Y' or 'N'
	PriceDesc             string // eg. 'Club';
	SeatsAvail            string // eg. '2';
	CategoryRange         string // eg. '';
	intCategoryMaxTickets string // eg. '2';
}
