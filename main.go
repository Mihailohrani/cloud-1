package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"cloud-1/handlers"
)

func main() {
	startTime := time.Now()

	mux := http.NewServeMux()

	// Status
	mux.HandleFunc("/countryinfo/v1/status/", handlers.StatusHandler(startTime))

	// Info
	mux.HandleFunc("/countryinfo/v1/info/", handlers.InfoHandler)

	// Exchange
	mux.HandleFunc("/countryinfo/v1/exchange/", handlers.ExchangeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
