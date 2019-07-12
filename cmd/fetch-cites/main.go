package main

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/core"
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
	pytmCities, err := ptm.FetchCities()
	if err != nil {
		panic(err)
	}
	fmt.Printf("PayTM cities: %+v\n", pytmCities)

	commonCities := core.MergeCities(bmsCities, pytmCities)
	fmt.Printf("commonCities: %+v\n", commonCities)
}
