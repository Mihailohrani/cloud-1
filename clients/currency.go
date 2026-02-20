package clients

import (
	"cloud-1/models"
	"cloud-1/utils"
)

// Base endpoint for the Currency API (expects ISO 4217 currency code at the end).
const currencyBase = "http://129.241.150.113:9090/currency/"

// GetRates retrieves the latest exchange rates using the given base currency.
// It returns a map where the key is the target currency code and the value is the rate.
// The HTTP status from the upstream service is returned for error handling in handlers.
func GetRates(base string) (map[string]float64, int, error) {
	var resp models.CurrencyResponse

	status, err := utils.GetJSON(currencyBase+base, &resp)
	if err != nil {
		return nil, status, err
	}

	return resp.Rates, status, nil
}
