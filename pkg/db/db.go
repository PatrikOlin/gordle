package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host     = "192.168.1.252"
	port     = 5432
	user     = "postgres"
	password = "bokhylla"
	dbname   = "ordle"
)

var DBClient *sqlx.DB

func Open() (*sqlx.DB, error) {
	var err error
	DBClient, err = sqlx.Connect("postgres", getPqslInfo())

	if err != nil {
		fmt.Println("ingen db")
		log.Fatalln(err)
	}

	return DBClient, nil
}

func getPqslInfo() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env")
	}

	return os.ExpandEnv("host=${DB_HOST} user=${DB_USER} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable")

}
