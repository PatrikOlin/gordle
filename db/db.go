package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DBClient *sqlx.DB

func Open() (*sqlx.DB, error) {
	DBClient, err := sqlx.Connect("sqlite3", "./db/data.db")

	if err != nil {
		fmt.Println("ingen db")
		log.Fatalln(err)
	}

	return DBClient, nil
}
