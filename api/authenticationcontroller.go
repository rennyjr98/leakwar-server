package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"leakwarsvr/persistence"
	"leakwarsvr/persistence/models"
)

var tokenizer = persistence.Tokenizer{
	LocalUrl:  "http://localhost:3000/api/token",
	OnlineUrl: "",
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/callback?approval_prompt=force",
	ClientID:     "311655574441-82u7t2emmn32c1s7mv2df7psu51040pk.apps.googleusercontent.com",
	ClientSecret: "1gO12qpPTIEi7YpIhr6Ed3cf",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/photoslibrary.readonly",
	},
	Endpoint: google.Endpoint,
}

// GoogleAuth : Method to start the autentication
func GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var authCode = r.FormValue("authCode")
	var user = r.FormValue("user")
	var token = getTokenFromWeb(authCode)
	tokenizer.Write(token, user)
	jsonToken, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(jsonToken))
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(authCode string) *oauth2.Token {
	tok, err := googleOauthConfig.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// GooglePhoto : Get a random photo from cloud
func GooglePhoto(w http.ResponseWriter, r *http.Request) {
	var jsonToken = r.FormValue("token")
	var token oauth2.Token
	err := json.Unmarshal([]byte(jsonToken), &token)
	if err != nil {
		log.Fatal(err)
	}

	var client = googleOauthConfig.Client(context.Background(), &token)
	response, err := client.Get("https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=100")
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var page models.GPage
	json.Unmarshal(content, &page)
	if len(page.MediaItems) > 0 {
		var index int = rand.Intn(len(page.MediaItems))
		var selectedPhoto = page.MediaItems[index]
		fmt.Println(string(selectedPhoto.BaseURL))
		fmt.Fprintf(w, string(selectedPhoto.BaseURL))
	} else {
		fmt.Fprintf(w, "")
	}
}
