package subscriptions

import "github.com/jaydp17/movie-ticket-watcher/pkg/db"

func WriteOne(s Subscription) error {
	subscriptions := make([]Subscription, 1)
	subscriptions[0] = s
	return Write(subscriptions)
}

func Write(subscriptions []Subscription) error {
	writables := make([]db.Writable, len(subscriptions))
	for i, s := range subscriptions {
		writables[i] = s
	}
	return db.Write(writables, TableName)
}
