package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jaydp17/movie-ticket-watcher/pkg/httperror"
)

type Payload struct {
	CityID   string `json:"cityID"`
	CinemaID string `json:"cinemaID"`
	MovieID  string `json:"movieID"`
	Date     string `json:"date"`
}

func validate(req events.APIGatewayProxyRequest) (Payload, error) {
	if len(req.Body) == 0 {
		return Payload{}, httperror.New(400, "body can't be empty")
	}

	var payload Payload
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return Payload{}, err
	}

	if len(payload.CityID) == 0 {
		return Payload{}, httperror.New(400, "cityID can't be empty")
	}
	if len(payload.CinemaID) == 0 {
		return Payload{}, httperror.New(400, "cinemaID can't be empty")
	}
	if len(payload.MovieID) == 0 {
		return Payload{}, httperror.New(400, "movieID can't be empty")
	}
	if len(payload.Date) == 0 {
		return Payload{}, httperror.New(400, "date can't be empty")
	}
	return payload, nil
}
