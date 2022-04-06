package daily

import (
	"log"
	"time"

	"github.com/PatrikOlin/gordle/pkg/db"
	"github.com/PatrikOlin/gordle/pkg/word"
)

type Daily struct {
	ID   int       `json:"id" db:"id"`
	Word string    `json:"word" db:"word"`
	Date time.Time `json:"date" db:"date"`
}

func GetDailyWord() (string, error) {
	var word string

	err := db.DBClient.Get(&word, "SELECT word FROM daily_word WHERE date = $1 ORDER BY id DESC LIMIT 1", time.Now().UTC().Format("01-02-2006"))
	if err != nil {
		return word, err
	}

	return word, nil
}

func Create() {
	d := Daily{
		Word: word.New(),
		Date: time.Now().UTC(),
	}

	persistDailyWord(d)
}

func persistDailyWord(d Daily) {
	stmt := "INSERT INTO daily_word (word, date) VALUES ($1, $2)"

	_, err := db.DBClient.Exec(stmt, d.Word, d.Date.Format("01-02-2006"))
	if err != nil {
		log.Fatalln(err)
	}
}
