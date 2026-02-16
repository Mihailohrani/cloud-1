package models

type StatusResponse struct {
	RestCountries int    `json:"restcountriesapi"`
	Currencies    int    `json:"currenciesapi"`
	Version       string `json:"version"`
	Uptime        int64  `json:"uptime"`
}

type CountryInfoResponse struct {
	Name       string            `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Area       float64           `json:"area"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    string            `json:"capital"`
}

type ExchangeResponse struct {
	Country       string               `json:"country"`
	BaseCurrency  string               `json:"base-currency"`
	ExchangeRates []map[string]float64 `json:"exchange-rates"`
}
