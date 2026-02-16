package main

import (
	"log"
	"net/http"
	"time"

	"cloud-1/handlers"
)

var startTime = time.Now()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/countryinfo/v1/status/", handlers.StatusHandler)
	mux.HandleFunc("/countryinfo/v1/info/", handlers.InfoHandler)
	mux.HandleFunc("/countryinfo/v1/exchange/", handlers.ExchangeHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
