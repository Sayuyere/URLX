package main

import (
	"math/rand"
	"os"
	"time"

	"urlx/api"
	"urlx/logging"
	"urlx/shortener"
	"urlx/store"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logger := logging.NewLogger()
	var s store.Store
	ps, err := store.NewPostgresStore()
	if err != nil {
		logger.Error("failed to connect to postgres: %v", err)
		os.Exit(1)
	}
	s = ps
	var shortenerSvc shortener.Shortener = shortener.NewSimpleShortener()

	r := api.SetupRouter(s, shortenerSvc, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Listening on :%s", port)
	r.Run(":" + port)
}
