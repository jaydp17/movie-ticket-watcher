package notifications

import (
	"encoding/json"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
	"os"
)

type pushPayload struct {
	Title   string       `json:"title,omitempty"`
	Body    string       `json:"body,omitempty"`
	Icon    string       `json:"icon,omitempty"`
	Vibrate []int        `json:"vibrate,omitempty"`
	Actions []pushAction `json:"actions,omitempty"`
}
type pushAction struct {
	Action string `json:"action,omitempty"`
	Title  string `json:"title,omitempty"`
	Icon   string `json:"icon,omitempty"`
}

func WebPush(result subscriptions.AvailableTicketResult) error {
	s := &webpush.Subscription{}
	if err := json.Unmarshal([]byte(result.Subscription.WebPushSubscription), s); err != nil {
		return InvalidWebPushSubscriptionError{webpushSubscription: result.Subscription.WebPushSubscription, Err: err}
	}

	vapidPublicKey, vapidPrivateKey := getVapidKeys()

	// Send Notification
	payload := pushPayload{
		Title:   fmt.Sprintf("%s: tickets available 🎉", result.Movie.Title),
		Body:    fmt.Sprintf("The tickets are available in %s", result.Cinema.Name),
		Vibrate: []int{200, 100, 200, 100, 200, 100, 400},
	}
	bytes, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return marshalErr
	}
	_, err := webpush.SendNotification(bytes, s, &webpush.Options{
		Subscriber:      "example@example.com",
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", payload)
	return nil
}

func getVapidKeys() (string, string) {
	publicKey := os.Getenv("VAPID_PUBLIC_KEY")
	privateKey := os.Getenv("VAPID_PRIVATE_KEY")

	if len(publicKey) == 0 {
		panic("env var VAPID_PUBLIC_KEY not passed")
	}
	if len(privateKey) == 0 {
		panic("env var VAPID_PRIVATE_KEY not passed")
	}
	return publicKey, privateKey
}

type InvalidWebPushSubscriptionError struct {
	webpushSubscription string
	Err                 error
}

func (e InvalidWebPushSubscriptionError) Error() string {
	return fmt.Sprintf("Got (%s), which isn't a valid webpush subscription object, which lead to Error: %+v", e.webpushSubscription, e.Err)
}
