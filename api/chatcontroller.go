package api

import (
	"encoding/json"
	"fmt"
	"leakwarsvr/domain"

	socketio "github.com/googollee/go-socket.io"
)

// ChatController : Http access for ChatRoom entity
type ChatController struct {
	socket *socketio.Server
}

// Init : Initialize the structure
func (controller *ChatController) Init(socket *socketio.Server) {
	controller.socket = socket
	controller.socket.OnEvent("/", "msg", controller.SendMessage)
	fmt.Println("Hola mundo")
}

// SendMessage : Send message to room
func (controller *ChatController) SendMessage(s socketio.Conn, chatcast string) {
	var chat domain.ChatCast
	json.Unmarshal([]byte(chatcast), &chat)
	fmt.Println("I recieved: ", chat.Msg)
	controller.socket.BroadcastToRoom("", chat.To, "msg", chat.Msg)
}
