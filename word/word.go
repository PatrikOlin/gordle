package word

import (
	"log"

	"github.com/PatrikOlin/gordle/db"
)

func New() string {
	var word string

	err := db.DBClient.Get(&word, "SELECT word FROM words ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		log.Fatalln(err)
	}

	return word
}
