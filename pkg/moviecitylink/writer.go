package moviecitylink

import "github.com/jaydp17/movie-ticket-watcher/pkg/db"

func WriteOne(link MovieCityLink) error {
	links := make([]MovieCityLink, 0)
	links = append(links, link)
	return Write(links)
}

func Write(links []MovieCityLink) error {
	writables := make([]db.Writable, len(links))
	for i, l := range links {
		writables[i] = l
	}
	return db.Write(writables, TableName)
}
