package main

import "github.com/jaydp17/movie-ticket-watcher/pkg/providers"

func main() {
	bms := providers.BookMyShow{}
	bms.GetAvailableVenueCodes("ET00106002", "BANG", "20190705")
}
