package movies

import "github.com/jaydp17/movie-ticket-watcher/pkg/db"

func Write(movies []Movie) error {
	writables := make([]db.Writable, len(movies))
	for i, m := range movies {
		writables[i] = m
	}
	return db.Write(writables, TableName)
}
