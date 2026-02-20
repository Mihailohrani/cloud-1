package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"cloud-1/clients"
	"cloud-1/models"
)

// ExchangeHandler returns exchange rates from the base country currency
// to the currencies of its neighboring countries.
func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/exchange/")
	code = strings.Trim(code, "/")
	code = strings.ToLower(code)

	if code == "" {
		http.Error(w, "Missing country code. Use /countryinfo/v1/exchange/{two_letter_code}", http.StatusBadRequest)
		return
	}

	// 1) Fetch base country data
	base, status, err := clients.GetCountryByAlpha(code)
	if err != nil {
		if status == http.StatusNotFound {
			http.Error(w, "Country not found", http.StatusNotFound)
			return
		}
		http.Error(w, "REST Countries API error", http.StatusBadGateway)
		return
	}

	baseCurrency := firstCurrencyCode(base.Currencies)
	if baseCurrency == "" {
		http.Error(w, "No base currency available for country", http.StatusBadGateway)
		return
	}

	// No borders → no exchange rates
	if len(base.Borders) == 0 {
		out := models.ExchangeResponse{
			Country:       base.Name.Common,
			BaseCurrency:  baseCurrency,
			ExchangeRates: []map[string]float64{},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
		return
	}

	// 2) Collect unique neighbour currencies
	neighborCurrenciesSet := map[string]struct{}{}
	for _, alpha3 := range base.Borders {
		alpha3 = strings.TrimSpace(alpha3)
		if alpha3 == "" {
			continue
		}

		neighbor, nStatus, err := clients.GetCountryByAlpha(alpha3)
		if err != nil {
			if nStatus == http.StatusNotFound {
				continue
			}
			http.Error(w, "REST Countries API error", http.StatusBadGateway)
			return
		}

		ccy := firstCurrencyCode(neighbor.Currencies)
		if ccy != "" && ccy != baseCurrency {
			neighborCurrenciesSet[ccy] = struct{}{}
		}
	}

	neighborCurrencies := make([]string, 0, len(neighborCurrenciesSet))
	for c := range neighborCurrenciesSet {
		neighborCurrencies = append(neighborCurrencies, c)
	}
	sort.Strings(neighborCurrencies)

	// 3) Retrieve exchange rates for the base currency
	rates, _, err := clients.GetRates(baseCurrency)
	if err != nil {
		http.Error(w, "Currencies API error", http.StatusBadGateway)
		return
	}

	// Build response in required format
	ratesOut := make([]map[string]float64, 0, len(neighborCurrencies))
	for _, ccy := range neighborCurrencies {
		if rate, ok := rates[ccy]; ok {
			ratesOut = append(ratesOut, map[string]float64{ccy: rate})
		}
	}

	out := models.ExchangeResponse{
		Country:       base.Name.Common,
		BaseCurrency:  baseCurrency,
		ExchangeRates: ratesOut,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// firstCurrencyCode returns the first currency code found in the map.
func firstCurrencyCode(m map[string]struct{}) string {
	for k := range m {
		return k
	}
	return ""
}
