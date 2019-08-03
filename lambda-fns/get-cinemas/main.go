package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
)

type Response = events.APIGatewayProxyResponse

func Handler(cityID string) ([]cinemas.Cinema, error) {
	outputCh := cities.FindByID(context.Background(), cityID)
	result := <-outputCh
	if result.Err != nil {
		fmt.Printf("error fetching city: %+v", result.Err)
		return nil, httperror.New(404, "can't find that city")
	}
	cinemasInTheCity := cinemas.Fetch(result.City)
	return cinemasInTheCity, nil
}

func HandlerBoilerplate(req events.APIGatewayProxyRequest) (Response, error) {
	err := validate(req)
	if err != nil {
		return lambdautils.ToResponse(err)
	}

	cityID := req.PathParameters["cityID"]
	result, err := Handler(cityID)
	if err != nil {
		return lambdautils.ToResponse(err)
	}
	return lambdautils.ToResponse(result)
}

func validate(req events.APIGatewayProxyRequest) error {
	cityID, ok := req.PathParameters["cityID"]
	if !ok || len(cityID) == 0 {
		return httperror.New(400, "cityID is required")
	}
	return nil
}

func main() {
	lambda.Start(HandlerBoilerplate)
}
