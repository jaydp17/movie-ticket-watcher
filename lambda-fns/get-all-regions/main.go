package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jaydp17/movie-ticket-watcher/pkg/cities"
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

type Response events.APIGatewayProxyResponse

// HandlerBoilerplate is our lambda handler invoked by the `lambda.Start` function call
func HandlerBoilerplate(_ context.Context) (Response, error) {
	result := Handler()

	body, err := json.Marshal(result)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers:         map[string]string{"Content-Type": "application/json"},
	}

	return resp, nil
}

func main() {
	lambda.Start(HandlerBoilerplate)
}
