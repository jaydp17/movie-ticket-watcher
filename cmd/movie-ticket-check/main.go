package main

import (
	"fmt"
	"log"

	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
)

func main() {
	bms := bookmyshow.BookMyShow{}
	showDetails, err := bms.GetAvailableVenueCodes("ET00106002", "BANG", "20190705")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", showDetails)
}
