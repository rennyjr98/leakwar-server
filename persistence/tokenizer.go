package persistence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// Tokenizer : Token file manager
type Tokenizer struct {
	LocalUrl  string
	OnlineUrl string
}

type IDToken struct {
	User  string
	Token *oauth2.Token
}

func (manager *Tokenizer) Write(token *oauth2.Token, user string) {
	port := os.Getenv("PORT")
	url := ""
	if port == "" {
		url = manager.LocalUrl
	} else {
		url = manager.OnlineUrl
	}

	idToken := IDToken{
		User:  user,
		Token: token,
	}
	reqBody, err := json.Marshal(idToken)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp)
	}
}
