package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
	"github.com/jaydp17/movie-ticket-watcher/pkg/logger"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
)

type Response = events.APIGatewayProxyResponse

func Handler(cityID string) ([]movies.Movie, error) {
	dbClient := db.NewClient()
	outputCh := cities.FindByID(context.Background(), dbClient, cityID)
	result := <-outputCh
	if result.Err != nil {
		log := logger.New()
		log.Errorf("error fetching city: %+v\n", result.Err)
		return nil, httperror.New(404, "can't find that city")
	}
	moviesInTheCity := movies.Fetch(dbClient, result.City)
	return moviesInTheCity, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func HandlerBoilerplate(req events.APIGatewayProxyRequest) (Response, error) {
	err := validate(req)
	if err != nil {
		return lambdautils.ToResponse(err)
	}

	regionCode := req.PathParameters["regionCode"]
	result, err := Handler(regionCode)
	if err != nil {
		return lambdautils.ToResponse(err)
	}
	return lambdautils.ToResponse(result)
}

func validate(req events.APIGatewayProxyRequest) error {
	regionCode, ok := req.PathParameters["regionCode"]
	if !ok || len(regionCode) == 0 {
		return fmt.Errorf("regionCode is required")
	}
	return nil
}

func main() {
	lambda.Start(HandlerBoilerplate)
}
