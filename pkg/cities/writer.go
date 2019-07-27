package cities

import (
	"github.com/jaydp17/movie-ticket-watcher/pkg/db"
)

func Write(cities []City) error {
	writables := make([]db.Writable, len(cities))
	for i, city := range cities {
		writables[i] = city
	}
	return db.Write(writables, TableName)
}
