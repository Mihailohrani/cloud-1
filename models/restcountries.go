package models

type RestCountry struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`

	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Area       float64           `json:"area"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`

	Flags struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"flags"`

	Capital []string `json:"capital"`

	Currencies map[string]struct{} `json:"currencies"`
}
