package providers

import (
	"fmt"
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
	"github.com/jaydp17/movie-ticket-watcher/pkg/utils"
	"strings"
)

// Provider is composed of Multiple single function interfaces to make it easier to mock while testing
type Provider interface {
	CitiesFetcher
	MoviesAndCinemasFetcher
	AvailableVenueCodesFetcher
}

type CitiesFetcher interface {
	FetchCities() ([]City, error)
}

type MoviesAndCinemasFetcher interface {
	FetchMoviesAndCinemas(city City) ([]Movie, []Cinema, error)
}

type AvailableVenueCodesFetcher interface {
	FetchAvailableVenueCodes(providerCityID, providerMovieID string, date db.YYYYMMDDTime) <-chan VenueCodesResult
}

type City struct {
	ID        string
	Name      string
	IsTopCity bool
}

type Cinema struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	CityID    string  `json:"cityID"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

func (c Cinema) NameSlug() string {
	return strings.ToLower(utils.KeepJustAlphaNumeric(c.Name))
}

type Movie struct {
	ID           string `json:"id"`
	GroupID      string `json:"groupID"`
	Title        string `json:"title"` // make sure there's no additional information in the title like (3D) or the Language
	ScreenFormat string `json:"screenFormat"`
	Language     string `json:"language"`
	ImageURL     string `json:"imageURL"`
}

// Slug combines the title-screenformat-language and returns that as a slug
func (m Movie) Slug() string {
	normalizedTitle := m.GroupSlug()
	screenFormat := utils.KeepJustAlphaNumeric(m.ScreenFormat)
	lang := utils.KeepJustAlphaNumeric(m.Language)
	slug := strings.ToLower(fmt.Sprintf("%s-%s-%s", normalizedTitle, screenFormat, lang))
	return slug
}

// GroupSlug is nothing but a normalized title without any screen format or language
// It's just the name of the movie normalized
func (m Movie) GroupSlug() string {
	cleanTitle := utils.KeepJustAlphaNumeric(m.Title)
	return strings.ToLower(cleanTitle)
}

type VenueCodesResult struct {
	Data []string
	Err  error
}
