package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cinemas"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
	"github.com/jaydp17/movie-ticket-watcher/pkg/logger"
	"sort"
)

type Response = events.APIGatewayProxyResponse

func Handler(cityID string) ([]cinemas.Cinema, error) {
	dbClient := db.NewClient()
	outputCh := cities.FindByID(context.Background(), dbClient, cityID)
	result := <-outputCh
	if result.Err != nil {
		log := logger.New()
		log.Errorf("error fetching city: %+v\n", result.Err)
		return nil, httperror.New(404, "can't find that city")
	}
	cinemasInTheCity := cinemas.Fetch(dbClient, result.City)
	sort.Slice(cinemasInTheCity, func(i, j int) bool {
		return cinemasInTheCity[i].ID < cinemasInTheCity[j].ID
	})
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
