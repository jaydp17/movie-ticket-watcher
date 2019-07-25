package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/lambdautils"
	"github.com/jaydp17/movie-ticket-watcher/pkg/movies"
)

type Response = events.APIGatewayProxyResponse

func Handler(cityID string) []movies.Movie {
	return make([]movies.Movie, 0)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func HandlerBoilerplate(req events.APIGatewayProxyRequest) (Response, error) {
	err := validate(req)
	if err != nil {
		return lambdautils.ToResponse(err)
	}

	regionCode := req.PathParameters["regionCode"]
	result := Handler(regionCode)
	return lambdautils.ToResponse(result)
}

func validate(req events.APIGatewayProxyRequest) error {
	regionCode, ok := req.PathParameters["regionCode"]
	if !ok || len(regionCode) == 0 {
		return fmt.Errorf("regionCode is required")
	}
	return fmt.Errorf("dummy error")
}

func main() {
	lambda.Start(HandlerBoilerplate)
}
