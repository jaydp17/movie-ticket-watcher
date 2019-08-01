package main

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
)

func main() {
	bms := bookmyshow.Provider{}
	bmsCities, err := bms.FetchCities()
	if err != nil {
		panic(err)
	}
	fmt.Printf("BMS cities: %+v\n", bmsCities)

	ptm := paytm.Provider{}
	ptmCities, err := ptm.FetchCities()
	if err != nil {
		panic(err)
	}
	fmt.Printf("PayTM cities: %+v\n", ptmCities)

	commonCities := cities.Merge(bmsCities, ptmCities)
	if err := cities.Write(commonCities); err != nil {
		panic(err)
	}
}
