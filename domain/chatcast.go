package domain

import (
	"leakwarsvr/persistence/models"
)

//ChatCast : Structure for send mensages
type ChatCast struct {
	Author models.User `json:"author"`
	Msg    string      `json:"msg"`
	To     string      `json:"to"`
}
