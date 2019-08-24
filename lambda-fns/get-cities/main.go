package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
	"sort"
)

// Handler has the core logic
func Handler() []cities.City {
	dbClient := db.NewClient()
	citiesChan := cities.All(dbClient)
	allCities := make([]cities.City, 0)
	for city := range citiesChan {
		allCities = append(allCities, city)
	}
	sort.Slice(allCities, func(i, j int) bool {
		return allCities[i].ID < allCities[j].ID
	})
	return allCities
}

type Response = events.APIGatewayProxyResponse

// HandlerBoilerplate is our lambda handler invoked by the `lambda.Start` function call
func HandlerBoilerplate(_ context.Context) (Response, error) {
	result := Handler()
	return lambdautils.ToResponse(result)
}

func main() {
	lambda.Start(HandlerBoilerplate)
}
