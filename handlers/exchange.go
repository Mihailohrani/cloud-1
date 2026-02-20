package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"cloud-1/models"
)

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

	client := &http.Client{Timeout: 8 * time.Second}

	// 1) Base country
	baseURL := "http://129.241.150.113:8080/v3.1/alpha/" + code
	baseResp, err := client.Get(baseURL)
	if err != nil {
		http.Error(w, "Failed to call REST Countries API", http.StatusBadGateway)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(baseResp.Body)

	if baseResp.StatusCode != http.StatusOK {
		if baseResp.StatusCode == http.StatusNotFound {
			http.Error(w, "Country not found", http.StatusNotFound)
			return
		}
		http.Error(w, "REST Countries API error", http.StatusBadGateway)
		return
	}

	var baseCountries []models.RestCountry
	if err := json.NewDecoder(baseResp.Body).Decode(&baseCountries); err != nil || len(baseCountries) == 0 {
		http.Error(w, "Failed to decode country data", http.StatusBadGateway)
		return
	}
	base := baseCountries[0]

	baseCurrency := firstCurrencyCode(base.Currencies)
	if baseCurrency == "" {
		http.Error(w, "No base currency available for country", http.StatusBadGateway)
		return
	}

	// No borders, return empty exchange rates
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

	// 2) Neighbor currencies
	neighborCurrenciesSet := map[string]struct{}{}
	for _, alpha3 := range base.Borders {
		alpha3 = strings.TrimSpace(alpha3)
		if alpha3 == "" {
			continue
		}

		neighborURL := "http://129.241.150.113:8080/v3.1/alpha/" + alpha3
		nResp, err := client.Get(neighborURL)
		if err != nil {
			http.Error(w, "Failed to call REST Countries API", http.StatusBadGateway)
			return
		}

		if nResp.StatusCode != http.StatusOK {
			_ = nResp.Body.Close()
			http.Error(w, "REST Countries API error", http.StatusBadGateway)
			return
		}

		var nCountries []models.RestCountry
		if err := json.NewDecoder(nResp.Body).Decode(&nCountries); err != nil || len(nCountries) == 0 {
			_ = nResp.Body.Close()
			http.Error(w, "Failed to decode neighbor country data", http.StatusBadGateway)
			return
		}
		_ = nResp.Body.Close()

		ccy := firstCurrencyCode(nCountries[0].Currencies)
		if ccy != "" && ccy != baseCurrency {
			neighborCurrenciesSet[ccy] = struct{}{}
		}
	}

	neighborCurrencies := make([]string, 0, len(neighborCurrenciesSet))
	for c := range neighborCurrenciesSet {
		neighborCurrencies = append(neighborCurrencies, c)
	}
	sort.Strings(neighborCurrencies)

	// 3) Currency API call (once)
	rates, err := fetchRates(client, baseCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Build output in the allowed format: []map[string]float64
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

// Picks the first key in the currencies map
func firstCurrencyCode(m map[string]struct{}) string {
	for k := range m {
		return k
	}
	return ""
}

func fetchRates(client *http.Client, base string) (map[string]float64, error) {
	// Correct endpoint per docs: /currency/{BASE}
	url := "http://129.241.150.113:9090/currency/" + base

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call Currencies API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return nil, fmt.Errorf("currencies API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	// Expected shape: {"base":"NOK","rates":{...}}
	var cr models.CurrencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err == nil && len(cr.Rates) > 0 {
		return cr.Rates, nil
	}

	// Fallback if JSON differs
	resp2, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call Currencies API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp2.Body)

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("currencies API error (status %d)", resp2.StatusCode)
	}

	var raw map[string]any
	if err := json.NewDecoder(resp2.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode currency data")
	}

	if rAny, ok := raw["rates"]; ok {
		if rMap, ok := rAny.(map[string]any); ok {
			out := map[string]float64{}
			for k, v := range rMap {
				if f, ok := v.(float64); ok {
					out[k] = f
				}
			}
			if len(out) > 0 {
				return out, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to decode currency rates")
}
