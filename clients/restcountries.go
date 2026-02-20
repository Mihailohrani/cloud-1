package clients

import (
	"fmt"

	"cloud-1/models"
	"cloud-1/utils"
)

// Base endpoint for REST Countries alpha code lookup.
const restCountriesBase = "http://129.241.150.113:8080/v3.1/alpha/"

// GetCountryByAlpha retrieves a country by its ISO alpha code from the REST Countries API.
// It returns the first matching country, the upstream HTTP status, and an error if the call fails.
func GetCountryByAlpha(code string) (models.RestCountry, int, error) {
	var countries []models.RestCountry

	status, err := utils.GetJSON(restCountriesBase+code, &countries)
	if err != nil {
		return models.RestCountry{}, status, err
	}

	if len(countries) == 0 {
		return models.RestCountry{}, status, fmt.Errorf("empty response")
	}

	return countries[0], status, nil
}
