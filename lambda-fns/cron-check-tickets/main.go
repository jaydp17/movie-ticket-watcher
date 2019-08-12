package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/notifications"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
)

func Handler() {
	dbClient := db.NewClient()

	allSubscriptions := make([]subscriptions.Subscription, 0)
	for subscription := range subscriptions.All(dbClient) {
		allSubscriptions = append(allSubscriptions, subscription)
	}

	bmsProvider := bookmyshow.Provider{}
	ptmProvider := paytm.Provider{}
	availableTickets := subscriptions.CheckForAvailableTickets(dbClient, bmsProvider, ptmProvider, allSubscriptions)
	for sub := range availableTickets {
		go notifications.WebPush(sub)
	}
}

func main() {
	//Handler()
	lambda.Start(Handler)
}
