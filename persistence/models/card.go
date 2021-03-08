package models

// Card : Struct that contains card info
type Card struct {
	Name         string `json:"name"`
	Author       string `json:"author"`
	Description  string `json:"description"`
	Size         int    `json:"size"`
	TurnEffect   int    `json:"turnEffect"`
	ChangePlayer bool   `json:"changePlayer"`
	ChangeCard   bool   `json:"changeCard"`
	IsSlave      bool   `json:"isSlave"`
	ToBlock      bool   `json:"toBlock"`
	Target       string `json:"target"`
}
