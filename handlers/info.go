package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"cloud-1/clients"
	"cloud-1/models"
)

// InfoHandler returns general country information for a given ISO alpha-2 code.
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

	// Fetch country data from REST Countries API via client
	c, status, err := clients.GetCountryByAlpha(code)
	if err != nil {
		if status == http.StatusNotFound {
			http.Error(w, "Country not found", http.StatusNotFound)
			return
		}
		http.Error(w, "REST Countries API error", http.StatusBadGateway)
		return
	}

	// Use first capital if available
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
