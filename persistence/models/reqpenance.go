package models

// ReqPenance : Structure that indicates the penance for a player
type ReqPenance struct {
	PublicList map[string]int `json:"publicList"`
	Slaves     []*Slave       `json:"slaves"`
	BlackList  []string       `json:"blackList"`
}
