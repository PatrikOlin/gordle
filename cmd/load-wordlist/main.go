package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/PatrikOlin/gordle/pkg/db"
)

func init() {
	_, err := db.Open()
	if err != nil {
		log.Fatalln("Failed to connect to database")
	}
}

func process() {
	f, err := os.Open("dedup-svenska-ord")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	fmt.Println("file read")
	count := 0
	for scanner.Scan() {
		stmt := "INSERT INTO word (word) VALUES ($1)"
		_, err = db.DBClient.Exec(stmt, scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		count++
		if count%100 == 0 {
			fmt.Println("count ", count)
			fmt.Println("Inserted ", scanner.Text())
		}
	}
}

func main() {
	process()
}
