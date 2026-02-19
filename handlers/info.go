package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"cloud-1/models"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/info/")
	code = strings.Trim(code, "/")
	code = strings.ToLower(code)

	if code == "" {
		http.Error(w, "Missing country code. Use /countryinfo/v1/info/{two_letter_code}", http.StatusBadRequest)
		return
	}

	url := "http://129.241.150.113:8080/v3.1/alpha/" + code

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to call REST Countries API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// keep it simple and predictable
		if resp.StatusCode == http.StatusNotFound {
			http.Error(w, "Country not found", http.StatusNotFound)
			return
		}
		http.Error(w, "REST Countries API error", http.StatusBadGateway)
		return
	}

	var countries []models.RestCountry
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil || len(countries) == 0 {
		http.Error(w, "Failed to decode country data", http.StatusBadGateway)
		return
	}

	c := countries[0]

	capital := ""
	if len(c.Capital) > 0 {
		capital = c.Capital[0]
	}

	response := models.CountryInfoResponse{
		Name:       c.Name.Common,
		Continents: c.Continents,
		Population: c.Population,
		Area:       c.Area,
		Languages:  c.Languages,
		Borders:    c.Borders,
		Flag:       c.Flags.PNG,
		Capital:    capital,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
