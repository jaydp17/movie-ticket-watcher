package cinemas

import "github.com/jaydp17/movie-ticket-watcher/pkg/db"

func Write(cinemas []Cinema) error {
	writables := make([]db.Writable, len(cinemas))
	for i, c := range cinemas {
		writables[i] = c
	}
	return db.Write(writables, TableName)
}
