package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
)

// Handler has the core logic
func Handler() []cities.City {
	citiesChan := cities.All()
	allCities := make([]cities.City, 0)
	for city := range citiesChan {
		allCities = append(allCities, city)
	}
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
