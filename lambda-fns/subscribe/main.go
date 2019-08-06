package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
	"github.com/jaydp17/movie-ticket-watcher/pkg/subscriptions"
)

type Response = events.APIGatewayProxyResponse

func Handler(payload Payload) (subscriptions.Subscription, error) {
	subscription, err := subscriptions.New(payload.CityID, payload.MovieID, payload.CinemaID, payload.WebPushSubscription, payload.ScreeningDate.Time)
	if err != nil {
		return subscriptions.Subscription{}, err
	}
	if err := subscriptions.WriteOne(subscription); err != nil {
		return subscriptions.Subscription{}, err
	}
	return subscription, nil
}

func HandlerBoilerplate(req events.APIGatewayProxyRequest) (Response, error) {
	payload, validationErr := validate(req)
	if validationErr != nil {
		fmt.Printf("error validating request: %+v", validationErr)
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
