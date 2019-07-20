package providers

import (
	"fmt"
)

// Provider is an interface
type Provider interface {
	FetchCities() ([]City, error)
	FetchMoviesAndCinemas(city City) ([]Movie, []Cinema, error)
}

type City struct {
	ID        string
	Name      string
	IsTopCity bool
}

type Cinema struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Provider  string    `json:"provider"`
	CityID    string    `json:"cityID"`
	Latitude  Latitude  `json:"latitude"`
	Longitude Longitude `json:"longitude"`
	Address   string    `json:"address"`
}
type Latitude float64
type Longitude = Latitude

func (l Latitude) String() string {
	strWithExtraDecimal := fmt.Sprintf("%.4f", l)
	return strWithExtraDecimal[:len(strWithExtraDecimal)-1]
}

type Movie struct {
	ID           string `json:"id"`
	GroupID      string `json:"groupID"`
	Title        string `json:"title"`
	ScreenFormat string `json:"screenFormat"`
	Language     string `json:"language"`
	ImageURL     string `json:"imageURL"`
}
