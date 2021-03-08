package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
)

// Cruppier : Struct to organize player's cards
type Cruppier struct {
	Deck []*Card
}

// Load : Load cards from database
func (cruppier *Cruppier) Load() {
	file, err := ioutil.ReadFile("database.json")
	if err != nil {
		log.Fatal(err)
	}

	var cards []*Card
	err = json.Unmarshal(file, &cards)
	if err != nil {
		log.Fatal(err)
	}
	cruppier.Deck = cards
}

// Get : Return a set of cards
func (cruppier *Cruppier) Get(size int) []*Card {
	var bundle []*Card
	for i := 0; i < size; i++ {
		index := rand.Intn(len(cruppier.Deck))
		bundle = append(bundle, cruppier.Deck[index])
	}
	return bundle
}

// Check : Verify that user has the card
func (cruppier *Cruppier) Check(player *User, card *Card) bool {
	return true
}
