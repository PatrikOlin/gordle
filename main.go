package main

import (
	"log"
	"net/http"

	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"github.com/PatrikOlin/gordle/api"
	"github.com/PatrikOlin/gordle/db"
)

var (
	// version string
	addr string
)

func init() {
	pflag.StringVarP(&addr, "address", "a", ":4040", "the address for the api to listen on. Host and port separated by ':'")
	pflag.Parse()

	_, err := db.Open()
	if err != nil {
		log.Fatalln("Failed to connect to database")
	}
}

// func process() {
// 	f, err := os.Open("swe_wordlist_filtered")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)

// 	fmt.Println("file read")
// 	count := 0
// 	for scanner.Scan() {
// 		stmt := "INSERT INTO words (word) VALUES (?)"
// 		_, err = db.DBClient.Exec(stmt, scanner.Text())
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		count++
// 		if count%100 == 0 {
// 			fmt.Println("count ", count)
// 			fmt.Println("Inserted ", scanner.Text())
// 		}
// 	}

// }

func main() {
	// logger, _ := zap.NewProduction(zap.WithCaller(false))
	logger, err := newLogger()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	r := api.GetRouter(logger)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func newLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"gordle.log",
		"stdout",
	}
	return cfg.Build(zap.WithCaller(true))
}
