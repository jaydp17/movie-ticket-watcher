package notifications

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
)

func WebPush(subscription subscriptions.Subscription) {
	fmt.Printf("%s is available in %s 🎉\n", subscription.MovieID, subscription.CinemaID)
}
