package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DBClient *sqlx.DB

func Open() (*sqlx.DB, error) {
	var err error
	DBClient, err = sqlx.Connect("sqlite3", "_data.db")

	if err != nil {
		fmt.Println("ingen db")
		log.Fatalln(err)
	}

	return DBClient, nil
}
