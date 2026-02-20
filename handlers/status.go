package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"cloud-1/models"
)

const (
	restCountriesProbe = "http://129.241.150.113:8080/v3.1/alpha/no"
	currencyProbe      = "http://129.241.150.113:9090/currency/NOK"
)

// StatusHandler reports upstream service availability and this service uptime.
func StatusHandler(start time.Time) http.HandlerFunc {
	client := &http.Client{Timeout: 7 * time.Second}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		restCode := probeStatus(client, restCountriesProbe)
		currCode := probeStatus(client, currencyProbe)

		resp := models.StatusResponse{
			RestCountriesAPI: restCode,
			CurrenciesAPI:    currCode,
			Version:          "v1",
			Uptime:           int(time.Since(start).Seconds()),
		}

		status := http.StatusOK
		if restCode != http.StatusOK || currCode != http.StatusOK {
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// probeStatus returns the HTTP status code of a GET request to the given URL.
func probeStatus(client *http.Client, url string) int {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0
	}

	res, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	return res.StatusCode
}
