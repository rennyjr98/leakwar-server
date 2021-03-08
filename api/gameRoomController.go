package api

import (
	"encoding/json"
	"fmt"
	"leakwarsvr/domain"
	"leakwarsvr/persistence/models"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

// GameRoomController : Http access for GameRoom entity
type GameRoomController struct {
	socket   *socketio.Server
	cruppier *models.Cruppier
	Rooms    map[string]*models.GameRoom
}

// Init : Initialize the structure
func (controller *GameRoomController) Init(socket *socketio.Server) {
	go controller.timer()
	controller.socket = socket
	controller.Rooms = map[string]*models.GameRoom{}
	controller.cruppier = &models.Cruppier{}
	controller.cruppier.Load()

	socket.OnConnect("/", controller.Join)
	socket.OnEvent("/", "reqStartGame", controller.StartGame)
	socket.OnEvent("/", "vote", controller.Vote)
	socket.OnEvent("/", "reqPenance", controller.ReqPaymentSocket)
	socket.OnEvent("/", "applyJob", controller.ApplyJob)
	socket.OnEvent("/", "reqRouletteAction", controller.RouletteAction)
	socket.OnEvent("/", "send_penance", controller.SendPayment)
	socket.OnEvent("/", "set_card", controller.SetCard)
	socket.OnEvent("/", "disconnecting", controller.Leave)
}

func (controller *GameRoomController) timer() {
	for true {
		for key, room := range controller.Rooms {
			if room.GameStarted {
				if !room.IsInPause {
					room.Clock--
					controller.socket.BroadcastToRoom("", key, "clock", strconv.Itoa(room.Clock))
					if room.Clock == 0 {
						controller.Roulette(key)
						room.Clock = 11
					}
				} else if room.FailedVotes+room.UpVotes == room.UserLen() {
					room.IsInPause = false
					if room.FailedVotes > room.UpVotes {
						json, err := json.Marshal(room.VoteUser)
						if err != nil {
							log.Fatal(err)
						}
						controller.ReqPayment(string(json))
					}

					room.FailedVotes = 0
					room.UpVotes = 0
					room.Clock = 1
				}
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
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
	fmt.Println(user)
	friend, _ := json.Marshal(user)
	controller.socket.BroadcastToRoom("", room, "add_friend", string(friend))
	return nil
}

// distributeCards :
func (controller *GameRoomController) distributeCards(roomID string) {
	deck := controller.Rooms[roomID].DistributeCards(controller.cruppier)
	if deck != nil {
		jsonDeck, err := json.Marshal(deck)
		if err != nil {
			log.Fatal(err)
		}
		controller.socket.BroadcastToRoom("", roomID, "get_deck", string(jsonDeck))
	}
}

// asignAdmin : Choose a admin
func (controller *GameRoomController) asignAdmin(user *models.User, roomID string) {
	if controller.Rooms[roomID].UserLen() == 1 {
		user.Admin = true
	}
}

// StartGame ;
func (controller *GameRoomController) StartGame(s socketio.Conn, roomID string) {
	if !controller.Rooms[roomID].GameStarted {
		controller.Rooms[roomID].GameStarted = true
		controller.distributeCards(roomID)
		controller.Roulette(roomID)
	}
}

// ResumeGame ;
func (controller *GameRoomController) Vote(s socketio.Conn, pair string) {
	if len(pair) > 0 {
		pairs := strings.Split(pair, ",")
		if len(pairs) == 2 {
			vote := pairs[0]
			roomID := pairs[1]

			switch vote {
			case "Cumplido":
				controller.Rooms[roomID].UpVotes++
				break
			case "Fallido":
				controller.Rooms[roomID].FailedVotes++
				break
			}
		}
	}
}

// Roulette : Http access to roulette gameroom
func (controller *GameRoomController) Roulette(roomID string) {
	var roulette = controller.Rooms[roomID].Roulette()
	var roulettePackage, err = json.Marshal(roulette)
	if err != nil {
		log.Fatal(err)
	}
	controller.socket.BroadcastToRoom("", roomID, "get_roulette", string(roulettePackage))
	if controller.Rooms[roomID].Turn >= 20 {
		controller.socket.BroadcastToRoom("", roomID, "end", "")
		delete(controller.Rooms, roomID)
	}
}

// RouletteAction : Get a random action
func (controller *GameRoomController) RouletteAction(s socketio.Conn, user string) {
	probability := rand.Intn(100)
	stock := models.StockExchange{}
	var player models.User
	var action *models.StockExchange
	err := json.Unmarshal([]byte(user), &player)

	if err != nil {
		log.Fatal(err)
	} else {
		if probability < 40 {
			action = stock.GetJob()
		} else if probability < 56 {
			if player.Job != nil {
				action = stock.GetTax(player.Job.Economy[1])
			} else {
				action = stock.GetTax(player.PaySize)
			}
		} else {
			action = stock.GetSituation(controller.Rooms[player.Room].Level)
		}

		jsonAction, _ := json.Marshal(action)
		if probability > 55 {
			controller.Rooms[player.Room].IsInPause = true
			controller.Rooms[player.Room].VoteUser = &player
			controller.socket.BroadcastToRoom("", player.Room, "get_action", string(jsonAction))
		} else {
			s.Emit("get_action", string(jsonAction))
		}
	}
}

func (controller *GameRoomController) ApplyJob(s socketio.Conn, user string) {
	var player models.User
	err := json.Unmarshal([]byte(user), &player)
	if err != nil {
		log.Fatal(err)
	}

	controller.Rooms[player.Room].SetJob(player.Email, player.Job)
}

// SetCard : Add card to apply in requirement
func (controller *GameRoomController) SetCard(s socketio.Conn, ticket string) {
	var ticketCard models.TicketCard
	err := json.Unmarshal([]byte(ticket), &ticketCard)
	if err != nil {
		log.Fatal(err)
	}
	controller.Rooms[ticketCard.Room].SetCard(&ticketCard.Card)
	controller.socket.BroadcastToRoom("", ticketCard.Room, "applied_card", ticketCard.Card.Name)
}

func (controller *GameRoomController) ReqPaymentSocket(s socketio.Conn, user string) {
	controller.ReqPayment(user)
}

// ReqPayment : Get request payment for users
func (controller *GameRoomController) ReqPayment(user string) {
	var player models.User
	err := json.Unmarshal([]byte(user), &player)
	if err != nil {
		log.Fatal(err)
	}

	controller.Rooms[player.Room].SetTreat(player.Email, player.Treat)
	var req = controller.Rooms[player.Room].RequestPayment()
	jsonReqPayment, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	controller.socket.BroadcastToRoom("", player.Room, "req_payment", string(jsonReqPayment))
	controller.Rooms[player.Room].RemoveTreat(player.Email, player.Treat)
}

// SendPayment : Send payment to show in a room
func (controller *GameRoomController) SendPayment(s socketio.Conn, payment string) {
	var penance domain.Penance
	err := json.Unmarshal([]byte(payment), &penance)
	if err != nil {
		log.Fatal(err)
	}
	controller.socket.BroadcastToRoom("", penance.Group, "show_payment", payment)
}

// Leave : Get out a disconnect player
func (controller *GameRoomController) Leave(s socketio.Conn, roomID string) {
	controller.Rooms[roomID].RemovePlayer(s.ID())
	if controller.Rooms[roomID].UserLen() == 0 {
		delete(controller.Rooms, roomID)
	}
}
