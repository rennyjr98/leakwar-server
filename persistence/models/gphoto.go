package models

// GPhoto : Struct for JSON google photo item
type GPhoto struct {
	ID       string `json:"id"`
	BaseURL  string `json:"baseUrl"`
	FileName string `json:"filename"`
}
