package models

// User : Player properties
type User struct {
	Admin   bool    `json:"admin"`
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	PayTurn int     `json:"payturn"`
	Cards   []*Card `json:"cards"`
	Slaves  []*User `json:"slaves"`
}
