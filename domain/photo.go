package domain

// Photo : Struct that defines the content of data transfer
type Photo struct {
	URL     string   `json:"url"`
	To      []string `json:"to"`
	InCloud bool     `json:"inCloud"`
}
