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

func main() {
	logger, _ := zap.NewProduction(zap.WithCaller(false))
	defer func() {
		_ = logger.Sync()
	}()

	r := api.GetRouter(logger)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalln(err)
	}
}
