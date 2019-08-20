package subscriptions

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
	"github.com/jaydp17/movie-ticket-watcher/pkg/providers"
	"sync"
)

type AvailableTicketResult struct {
	Subscription Subscription
	Movie        movies.Movie
	Cinema       cinemas.Cinema
}

func CheckForAvailableTickets(dbClient dynamodbiface.ClientAPI, bms providers.AvailableVenueCodesFetcher, ptm providers.AvailableVenueCodesFetcher, allSubscriptions []Subscription) <-chan AvailableTicketResult {
	groupOfSubscriptions := groupSimilarSubscriptions(allSubscriptions)
	outputCh := make(chan AvailableTicketResult)
	wg := sync.WaitGroup{}
	wg.Add(len(groupOfSubscriptions))
	for _, similarSubscriptions := range groupOfSubscriptions {
		go func(similarSubscriptions groupSubscriptions) {
			defer wg.Done()
			sub := similarSubscriptions.subscriptions[0]
			city, movie, cinema, err := sub.GetMovieCityAndCinema(dbClient)
			if err != nil {
				fmt.Printf("error in GetMovieCityAndCinema: %v", err)
				return
			}
			areAvailable, err := AreTicketsAvailable(bms, ptm, city, movie, cinema, sub.ScreeningDate)
			if err != nil {
				fmt.Printf("error in AreTicketsAvailable: %v", err)
				return
			}
			if areAvailable {
				for _, availableSubscription := range similarSubscriptions.subscriptions {
					outputCh <- AvailableTicketResult{
						Subscription: availableSubscription,
						Movie:        movie,
						Cinema:       cinema,
					}
				}
			}
		}(similarSubscriptions)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()
	return outputCh
}

func groupSimilarSubscriptions(subscriptions []Subscription) []groupSubscriptions {
	groupedSubscriptions := make([]groupSubscriptions, 0)
	for _, sub := range subscriptions {
		foundSimilar := false
		for i := 0; i < len(groupedSubscriptions); i++ {
			// here we just compare it with the 1st subscription in the group
			// as if the 1st one passed other others should be the same
			if groupedSubscriptions[i].subscriptions[0].IsSimilar(sub) {
				groupedSubscriptions[i].subscriptions = append(groupedSubscriptions[i].subscriptions, sub)
				foundSimilar = true
				break
			}
		}
		if !foundSimilar {
			groupedElements := []Subscription{sub}
			groupedSubscriptions = append(groupedSubscriptions, groupSubscriptions{subscriptions: groupedElements})
		}
	}
	return groupedSubscriptions
}

type groupSubscriptions struct {
	subscriptions []Subscription
}
