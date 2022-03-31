package main

import (
	"fmt"
	"log"

	d "github.com/PatrikOlin/gordle/pkg/daily"
	"github.com/PatrikOlin/gordle/pkg/db"
)

func init() {
	_, err := db.Open()
	if err != nil {
		log.Fatalln("Failed to connect to database")
	}
}

func main() {
	setDailyWord()
	clearDailyUserSessions()

	word := d.GetDailyWord()

	fmt.Println(word)
}

func setDailyWord() {
	d.Create()
}

func clearDailyUserSessions() {
	stmt := "UPDATE user_session us SET finished_daily = FALSE"

	_, err := db.DBClient.Exec(stmt)
	if err != nil {
		log.Fatalln(err)
	}
}
