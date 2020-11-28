package api

import (
	"encoding/json"
	"leakwarsvr/domain"
	"leakwarsvr/persistence/models"
	"log"
	"net/url"

	socketio "github.com/googollee/go-socket.io"
)

// GameRoomController : Http access for GameRoom entity
type GameRoomController struct {
	socket *socketio.Server
	Rooms  map[string]*models.GameRoom
}

// Init : Initialize the structure
func (controller *GameRoomController) Init(socket *socketio.Server) {
	controller.socket = socket
	controller.Rooms = map[string]*models.GameRoom{}
	socket.OnConnect("/", controller.Join)
	socket.OnEvent("/", "disconnecting", controller.Leave)
}

// Join : Join a player to room
func (controller *GameRoomController) Join(s socketio.Conn) error {
	s.SetContext("")
	params, err := url.ParseQuery(s.URL().RawQuery)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var user models.User
	var room string = params["room"][0]
	json.Unmarshal([]byte(params["user"][0]), &user)
	s.Join(room)

	_, found := controller.Rooms[room]
	if !found {
		controller.Rooms[room] = &models.GameRoom{}
	}

	user.Cards = []*models.Card{}
	controller.Rooms[room].Join(&user, s.ID())
	friend, _ := json.Marshal(user)
	//controller.asignAdmin(&user, room)
	controller.socket.BroadcastToRoom("", room, "add_friend", string(friend))
	return nil
}

// GetPlayers : Get a page of players
func (controller *GameRoomController) GetPlayers(domain.PaginatedRequest) {

}

// asignAdmin : Choose a admin
func (controller *GameRoomController) asignAdmin(user *models.User, roomID string) {
	if controller.Rooms[roomID].UserLen() == 1 {
		user.Admin = true
	}
}

// Leave : Get out a disconnect player
func (controller *GameRoomController) Leave(s socketio.Conn, room string) {
	controller.Rooms[room].RemovePlayer(s.ID())
	if controller.Rooms[room].UserLen() == 0 {
		delete(controller.Rooms, room)
	}
}
