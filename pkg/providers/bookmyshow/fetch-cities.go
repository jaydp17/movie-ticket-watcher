package bookmyshow

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
)

func (p Provider) FetchCities() ([]providers.City, error) {
	params := req.Param{
		"cmd": "DEREGIONLIST",
		"f":   "json",
		"et":  "ALL",
		"t":   token,
		"ch":  "mobile",
	}
	headers := req.Header{
		"lang":       "en",
		"User-Agent": okHTTPUserAgent,
	}
	res, err := req.Get(p.urlToFetchCities, params, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch BMS cities: %v", err)
	}
	var jsonRes bmsRegionListResponse
	if err := res.ToJSON(&jsonRes); err != nil {
		return nil, fmt.Errorf("failed to parse response from BMS: %v", err)
	}

	var cities []providers.City
	for _, region := range jsonRes.BookMyShow.TopCities {
		cities = append(cities, providers.City{
			ID:        region.RegionCode,
			Name:      region.RegionName,
			IsTopCity: true,
		})
	}
	for _, region := range jsonRes.BookMyShow.OtherCities {
		cities = append(cities, providers.City{
			ID:        region.RegionCode,
			Name:      region.RegionName,
			IsTopCity: false,
		})
	}
	return cities, nil
}

type bmsRegionListResponse struct {
	BookMyShow struct {
		TopCities   []bmsRegion
		OtherCities []bmsRegion
	}
}

type bmsRegion struct {
	RegionCode         string
	RegionName         string
	Seq                string // numbers disguised as strings, eg, '1'
	Lat                string // eg. '19.0760'
	Long               string // eg. '72.8777'
	AllowSales         string // eg. "Y" or "N"
	isOlaEnabled       string // eg. "Y" or "N"
	Alias              []string
	SubRegions         []bmsSubRegion
	RegionSlug         string
	RegionCallCenterNo string
}

type bmsSubRegion struct {
	SubRegionCode string
	SubRegionName string
	Seq           string // numbers disguised as strings, eg, '1'
	Lat           string // eg. '19.0760'
	Long          string // eg. '72.8777'
	AllowSales    string // eg. "Y" OR "N"
	SubRegionSlug string
}
