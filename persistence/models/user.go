package models

// User : Player properties
type User struct {
	Admin   bool     `json:"admin"`
	ID      string   `json:"id"`
	Room    string   `json:"room"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	PayTurn int      `json:"payturn"`
	PaySize int      `json:"paySize"`
	Cards   []*Card  `json:"cards"`
	Owners  []string `json:"owners"`
	Job     *Job     `json:"job"`
	Treat   *Job     `json:"treat"`
}
