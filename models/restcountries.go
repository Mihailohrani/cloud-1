package models

type RestCountry struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
}
