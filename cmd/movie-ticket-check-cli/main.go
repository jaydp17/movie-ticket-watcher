package main

import (
	"fmt"
	"log"

	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/utils"
)

const (
	regionCode          = "BANG"
	movieChildEventCode = "ET00106002"
	dateStr             = "20190705"
)

var venuesOfChoice = map[string]string{
	"CXBL": "Central Spirit Mall",
	"PVBN": "Forum Mall, Koramangala",
}

func main() {
	bms := bookmyshow.BookMyShow{}
	availableVenueCodes, err := bms.GetAvailableVenueCodes(movieChildEventCode, regionCode, dateStr)
	if err != nil {
		log.Fatalln(err)
	}
	chosenVenueCodes := utils.MapKeys(venuesOfChoice)
	intersectingCodes := utils.Intersection(availableVenueCodes, chosenVenueCodes)

	var intersectingVenueNames []string
	for _, venueCode := range intersectingCodes {
		if venueName, ok := venuesOfChoice[venueCode]; ok {
			intersectingVenueNames = append(intersectingVenueNames, venueName)
		}
	}

	if len(intersectingVenueNames) > 0 {
		fmt.Printf("Movie is available in the following Cinemas 🎉\n")
		for _, venueName := range intersectingVenueNames {
			fmt.Printf("👉 %v\n", venueName)
		}
	}
}
