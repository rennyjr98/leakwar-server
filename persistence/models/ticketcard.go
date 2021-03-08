package models

type TicketCard struct {
	Room string `json:"room"`
	Card Card   `json:"card"`
}
