package models

import "fmt"

// GameRoom : Room for users (player) connected
type GameRoom struct {
	whiteList    map[string]*User
	blackList    map[string]*User
	inEffectCard []*Card
	turnPlayers  []*User
	turn         int
}

// Init : Initialize the GameRoom structure
func (room *GameRoom) init() {
	if room.whiteList == nil {
		room.whiteList = make(map[string]*User)
		room.blackList = make(map[string]*User)
		room.inEffectCard = []*Card{}
		room.turnPlayers = []*User{}
		room.turn = 0
	}
}

// UserLen : Get whiteList size
func (room *GameRoom) UserLen() int {
	return len(room.whiteList)
}

// Join : Join player to room
func (room *GameRoom) Join(user *User, ID string) {
	room.init()
	user.Admin = len(room.whiteList) == 0
	user.ID = ID
	room.whiteList[user.Email] = user
}

// RemovePlayer : Remove disconnected player from room
func (room *GameRoom) RemovePlayer(id string) {
	var email string = ""
	for key, user := range room.whiteList {
		fmt.Println(email, key)
		if user.ID == id {
			email = key
		}
	}

	fmt.Println("Remove Player : ", email)

	if _, ok := room.whiteList[email]; ok {
		delete(room.whiteList, email)
	}
	if _, ok := room.blackList[email]; ok {
		delete(room.blackList, email)
	}
}
