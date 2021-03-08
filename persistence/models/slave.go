package models

// Slave: Struct to apply quotes to slaves
type Slave struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Quote int    `json:"quote"`
}
