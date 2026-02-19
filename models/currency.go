package models

// Flexible target shape, we will also support other shapes in code.
type CurrencyResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}
