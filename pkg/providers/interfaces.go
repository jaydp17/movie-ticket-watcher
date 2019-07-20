package providers

// Provider is an interface
type Provider interface {
	FetchCities() ([]City, error)
	FetchMoviesAndCinemas(city City) ([]string, error)
}

type City struct {
	ID        string
	Name      string
	IsTopCity bool
}
