package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"urlx/api"
	"urlx/shortener"
	"urlx/store"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var s store.Store
	ps, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	s = ps
	var shortenerSvc shortener.Shortener = shortener.NewSimpleShortener()

	r := api.SetupRouter(s, shortenerSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	r.Run(":" + port)
}
