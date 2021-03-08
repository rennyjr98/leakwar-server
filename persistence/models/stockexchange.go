package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
)

var jobs []*Job
var situations map[string][]string

const (
	CLASSIC  = "classic"
	HOT      = "hot"
	EXTREME  = "extreme"
	WHATEVER = "whatever"
)

// StockExchange : Struct to communicate an action
type StockExchange struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	A    string `json:"a"`
	RA   []int  `json:"ra"`
	B    string `json:"b"`
	RB   []int  `json:"rb"`
}

// Job : Struct for job entity
type Job struct {
	Name    string `json:"name"`
	Economy []int  `json:"economy"`
}

// Load : Load jobs from database
func (stock *StockExchange) Load() {
	file, err := ioutil.ReadFile("jobs.json")
	fileSt, errSt := ioutil.ReadFile("situations.json")
	if err != nil && errSt != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &jobs)
	errSt = json.Unmarshal(fileSt, &situations)
	if err != nil && errSt != nil {
		log.Fatal(err)
	}
}

// GetJob : Get a random job
func (stock *StockExchange) GetJob() *StockExchange {
	if len(jobs) == 0 {
		stock.Load()
	}

	jobA := jobs[rand.Intn(len(jobs))]
	jobB := jobs[rand.Intn(len(jobs))]
	_stock := &StockExchange{
		Name: "Elige un empleo:",
		Type: 1,
		A:    jobA.Name,
		RA:   jobA.Economy,
		B:    jobB.Name,
		RB:   jobB.Economy,
	}
	return _stock
}

// GetTax : Generate a random taxes
func (stock *StockExchange) GetTax(publicQuote int) *StockExchange {
	privateQuote := rand.Intn(8) + 1
	tax := &StockExchange{
		Name: "Impuestos",
		Type: 2,
		A:    "Soborno al banquero" + string(privateQuote),
		RA:   []int{privateQuote},
		B:    "Cuota pública" + string(publicQuote),
		RB:   []int{publicQuote},
	}
	return tax
}

func (stock *StockExchange) GetSituation(level int) *StockExchange {
	var situation string
	if len(jobs) == 0 {
		stock.Load()
	}

	switch level {
	case 0:
		situation = situations[CLASSIC][rand.Intn(len(situations[CLASSIC]))]
		break
	case 1:
		situation = situations[HOT][rand.Intn(len(situations[HOT]))]
		break
	case 2:
		situation = situations[EXTREME][rand.Intn(len(situations[EXTREME]))]
		break
	case 3:
		situation = situations[WHATEVER][rand.Intn(len(situations[WHATEVER]))]
		break
	default:
		situation = situations[CLASSIC][rand.Intn(len(situations[CLASSIC]))]
		break
	}

	tax := &StockExchange{
		Name: situation,
		Type: 3,
		A:    "Cumplido",
		RA:   []int{},
		B:    "Fallido",
		RB:   []int{},
	}
	return tax
}
