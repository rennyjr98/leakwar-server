package models

import (
	"fmt"
	"math/rand"
)

// GameRoom : Room for users (player) connected
type GameRoom struct {
	whiteList      map[string]*User
	blackList      map[string]*User
	inEffectCard   []*Card
	turnPlayers    []*User
	Turn           int
	Clock          int
	GameStarted    bool
	IsRouletteTurn bool
	IsInPause      bool
	Level          int
	UpVotes        int
	FailedVotes    int
	VoteUser       *User
}

// Init : Initialize the GameRoom structure
func (room *GameRoom) init() {
	if room.whiteList == nil {
		room.whiteList = make(map[string]*User)
		room.blackList = make(map[string]*User)
		room.inEffectCard = []*Card{}
		room.turnPlayers = []*User{}
		room.Turn = 0
		room.Clock = 11
		room.IsRouletteTurn = true
		room.IsRouletteTurn = false
		room.Level = 1
	}
}

// UserLen : Get whiteList size
func (room *GameRoom) UserLen() int {
	return len(room.whiteList)
}

func (room *GameRoom) FindUser(email string) *User {
	if _, ok := room.whiteList[email]; ok {
		return room.whiteList[email]
	}
	if _, ok := room.blackList[email]; ok {
		return room.blackList[email]
	}
	return nil
}

// Join : Join player to room
func (room *GameRoom) Join(user *User, ID string) {
	if _, ok := room.whiteList[user.Email]; !ok {
		room.init()
		user.Admin = len(room.whiteList) == 0
		user.ID = ID
		room.whiteList[user.Email] = user
	}
}

// DistributeCards : Asign deck of cards to each player
func (room *GameRoom) DistributeCards(cruppier *Cruppier) map[string]*User {
	for _, user := range room.whiteList {
		user.Cards = cruppier.Get(20)
	}
	return room.whiteList
}

// Roulette : Select one random player
func (room *GameRoom) Roulette() *User {
	room.turnUp()
	var n int = rand.Intn(room.UserLen())
	var index int = 0
	var user *User
	for _, value := range room.whiteList {
		if index == n {
			user = value
		}
		index++
	}
	user.PaySize = 1
	room.turnPlayers = append(room.turnPlayers, user)
	return user
}

func (room *GameRoom) turnUp() {
	for index, user := range room.turnPlayers {
		user.PayTurn--
		if user.PayTurn <= 0 {
			room.turnPlayers = append(room.turnPlayers[:index], room.turnPlayers[index+1:]...)
		}
	}
}

func (room *GameRoom) RemoveTreat(email string, treat *Job) {
	if _, ok := room.whiteList[email]; ok {
		room.whiteList[email].Treat = nil
	}
	if _, ok := room.blackList[email]; ok {
		room.blackList[email].Treat = nil
	}

	index := -1
	for i, email := range room.whiteList[email].Owners {
		if email == room.findBnk().Email {
			index = i
		}
	}

	if index == -1 {
		room.whiteList[email].Owners = []string{}
	} else {
		room.whiteList[email].Owners = append(room.whiteList[email].Owners[:index], room.whiteList[email].Owners[index+1:]...)
	}

}

// RequestPayment : Send call for get the content to share
func (room *GameRoom) RequestPayment() *ReqPenance {
	room.applyCards()
	req := &ReqPenance{
		PublicList: make(map[string]int),
		Slaves:     []*Slave{},
		BlackList:  []string{},
	}

	for _, user := range room.turnPlayers {
		fmt.Println(user.Email, user.Job)
		var bnk *User
		if user.Treat != nil && user.Treat.Name == "altf4" {
			bnk = room.findBnk()
			user.Owners = append(user.Owners, bnk.Email)
		} else {
			if user.Job == nil {
				req.PublicList[user.Email] = user.PaySize
				fmt.Println("Job nil : ", user.PaySize)
			} else if user.PaySize == 1 {
				req.PublicList[user.Email] = user.Job.Economy[1]
				fmt.Println("User PaySize == 1", user.Job.Economy[1])
			} else {
				req.PublicList[user.Email] = user.Job.Economy[1] + user.PaySize
				fmt.Println("User PaySize > 1", user.Job.Economy[1]+user.PaySize)
			}
		}

		for _, owner := range user.Owners {
			slave := &Slave{
				Name:  user.Email,
				Owner: owner,
			}

			if owner == bnk.Email {
				slave.Quote = user.Treat.Economy[0]
			} else {
				slave.Quote = user.PaySize
			}
			req.Slaves = append(req.Slaves, slave)
		}
	}

	for key := range room.blackList {
		req.BlackList = append(req.BlackList, key)
	}

	return req
}

func (room *GameRoom) findBnk() *User {
	bnk, admin := room.findBnkByList(room.whiteList)
	if bnk != nil {
		return bnk
	}

	bnk, admin2 := room.findBnkByList(room.blackList)
	if bnk != nil {
		return bnk
	}

	if admin != nil {
		return admin
	}
	return admin2
}

func (room *GameRoom) findBnkByList(list map[string]*User) (*User, *User) {
	var bnk *User
	var admin *User
	for _, user := range room.whiteList {
		if user.Job != nil {
			if user.Job.Name == "Banquero" {
				bnk = user
			}
		} else if user.Admin {
			admin = user
		}
	}

	return bnk, admin
}

func (room *GameRoom) applyCards() {
	for ic, card := range room.inEffectCard {
		if card.ToBlock {
			room.blackList[card.Target] = room.whiteList[card.Target]
		} else if card.ChangePlayer {
			var index int = 0
			for i, user := range room.turnPlayers {
				if user.Email == card.Author {
					index = i
				}
			}
			room.turnPlayers = append(room.turnPlayers[:index], room.turnPlayers[index+1:]...)
			room.turnPlayers = append(room.turnPlayers, room.whiteList[card.Target])
		} else if card.Size > 0 {
			for _, user := range room.turnPlayers {
				user.PaySize += card.Size
			}
		}

		card.TurnEffect--
		room.turnUpCards(card, ic)
	}
}

func (room *GameRoom) turnUpCards(card *Card, ic int) {
	if card.TurnEffect <= 0 {
		room.inEffectCard = append(room.inEffectCard[:ic], room.inEffectCard[ic+1:]...)
	}
}

// SetCard : Adding card for apply effect after Roulette
func (room *GameRoom) SetCard(card *Card) {
	room.inEffectCard = append(room.inEffectCard, card)
}

// SetJob : Adding job to user
func (room *GameRoom) SetJob(email string, job *Job) {
	if _, ok := room.whiteList[email]; ok {
		room.whiteList[email].Job = job
		fmt.Println("Set job : ", email)
	} else if _, ok := room.blackList[email]; ok {
		room.blackList[email].Job = job
		fmt.Println("Set job : ", email)
	}
}

// SetTreat : Adding job to user
func (room *GameRoom) SetTreat(email string, job *Job) {
	if _, ok := room.whiteList[email]; ok {
		room.whiteList[email].Treat = job
	} else if _, ok := room.blackList[email]; ok {
		room.blackList[email].Treat = job
	}
}

// RemovePlayer : Remove disconnected player from room
func (room *GameRoom) RemovePlayer(id string) {
	var email string = ""
	for key, user := range room.whiteList {
		if user.ID == id {
			email = key
		}
	}

	if _, ok := room.whiteList[email]; ok {
		delete(room.whiteList, email)
	}
	if _, ok := room.blackList[email]; ok {
		delete(room.blackList, email)
	}
}
