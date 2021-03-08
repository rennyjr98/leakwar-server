package main

import (
	"leakwarsvr/api"
	"log"
	"net/http"
	"os"
)

func main() {
	socket := api.SocketConnection{}
	socket.CreateSocket()

	go socket.Server.Serve()
	defer socket.Server.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/socket.io/", socket.CorsMiddleware(socket.Server))
	http.Handle("/", http.FileServer(http.Dir("api/public")))
	http.HandleFunc("/auth", api.GoogleAuth)
	http.HandleFunc("/gphotos", api.GooglePhoto)
	log.Println("Serving at 127.0.0.1:" + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
