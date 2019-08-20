package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/notifications"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/bookmyshow"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers/paytm"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
	"sync"
)

func Handler() {
	dbClient := db.NewClient()

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
				fmt.Printf("error while sending push notification: %+v\n", err)
			}
			if err := subscriptions.MoveToArchive(dbClient, result.Subscription); err != nil {
				fmt.Printf("%+v\n", err)
			}
		}(result)
	}
	wg.Wait()
}

func main() {
	//Handler()
	lambda.Start(Handler)
}
