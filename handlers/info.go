package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"cloud-1/models"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {

	code := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/info/")
	if code == "" {
		http.Error(w, "Missing country code", http.StatusBadRequest)
		return
	}

	url := "http://129.241.150.113:8080/v3.1/alpha/" + code

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to call REST Countries API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	var countries []models.RestCountry

	err = json.NewDecoder(resp.Body).Decode(&countries)
	if err != nil || len(countries) == 0 {
		http.Error(w, "Failed to decode country data", http.StatusInternalServerError)
		return
	}

	response := models.CountryInfoResponse{
		Name: countries[0].Name.Common,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
