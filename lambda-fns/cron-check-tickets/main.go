package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/logger"
	"github.com/jaydp17/movie-ticket-watcher/pkg/notifications"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
	"sync"
)

func Handler() {
	dbClient := db.NewClient()
	log := logger.New()

	allSubscriptions := make([]subscriptions.Subscription, 0)
	for subscription := range subscriptions.All(dbClient) {
		allSubscriptions = append(allSubscriptions, subscription)
	}

	bmsProvider := bookmyshow.New()
	ptmProvider := paytm.New()
	availableTickets := subscriptions.CheckForAvailableTickets(dbClient, bmsProvider, ptmProvider, allSubscriptions)

	wg := sync.WaitGroup{}
	for result := range availableTickets {
		wg.Add(1)
		go func(result subscriptions.AvailableTicketResult) {
			defer wg.Done()
			if err := notifications.WebPush(result); err != nil {
				log.Errorf("error while sending push notification: %+v\n", err)
			}
			if err := subscriptions.MoveToArchive(dbClient, result.Subscription); err != nil {
				log.Errorf("%+v\n", err)
			}
		}(result)
	}
	wg.Wait()
}

func main() {
	//Handler()
	lambda.Start(Handler)
}
