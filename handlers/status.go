package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	restCountriesBase = "http://129.241.150.113:8080/v3.1/"
	currencyBase      = "http://129.241.150.113:9090/currency/"
)

type StatusResponse struct {
	RestCountriesAPI int    `json:"restcountriesapi"`
	CurrenciesAPI    int    `json:"currenciesapi"`
	Version          string `json:"version"`
	Uptime           int    `json:"uptime"`
}

func StatusHandler(start time.Time) http.HandlerFunc {
	client := &http.Client{
		Timeout: 7 * time.Second,
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		restCode := probeStatus(client, restCountriesBase+"alpha/no")
		currCode := probeStatus(client, currencyBase+"NOK")

		resp := StatusResponse{
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

func probeStatus(client *http.Client, url string) int {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0
	}

	res, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer res.Body.Close()

	return res.StatusCode
}
