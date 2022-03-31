package main

import (
	"log"
	"net/http"

	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"github.com/PatrikOlin/gordle/pkg/api"
	"github.com/PatrikOlin/gordle/pkg/db"
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
