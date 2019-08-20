package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
	"github.com/jaydp17/movie-ticket-watcher/pkg/logger"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
)

type Response = events.APIGatewayProxyResponse

func Handler(payload Payload) (subscriptions.Subscription, error) {
	dbClient := db.NewClient()
	subscription, err := subscriptions.New(dbClient, payload.CityID, payload.MovieID, payload.CinemaID, payload.WebPushSubscription, payload.ScreeningDate.Time)
	if err != nil {
		return subscriptions.Subscription{}, err
	}

	if err := subscriptions.WriteOne(dbClient, subscription); err != nil {
		return subscriptions.Subscription{}, err
	}
	return subscription, nil
}

func HandlerBoilerplate(req events.APIGatewayProxyRequest) (Response, error) {
	payload, validationErr := validate(req)
	if validationErr != nil {
		log := logger.New()
		log.Errorf("error validating request: %+v\n", validationErr)
		return lambdautils.ToResponse(validationErr)
	}
	result, err := Handler(payload)
	if err != nil {
		return lambdautils.ToResponse(err)
	}
	return lambdautils.ToResponse(result)
}

func main() {
	lambda.Start(HandlerBoilerplate)
}
