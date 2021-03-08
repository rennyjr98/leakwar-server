package domain

//ChatCast : Structure for send mensages
type ChatCast struct {
	Author string `json:"author"`
	Msg    string `json:"msg"`
	To     string `json:"to"`
}
