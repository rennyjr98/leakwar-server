package domain

// Penance : Contents the penance of a user
type Penance struct {
	Author string  `json:"author"`
	Album  []Photo `json:"album"`
	Group  string  `json:"group"`
}
