package models

// Deck : Struct for send list of cards
type Deck struct {
	Author string  `json:"author"`
	Cards  []*Card `json:"cards"`
}
