package api

import (
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type controller interface {
	Init(*socketio.Server)
}

// SocketConnection : Has server properties
type SocketConnection struct {
	Server      *socketio.Server
	Controllers []controller
}

// CorsMiddleware : Control all request setting and headers
func (socket *SocketConnection) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
		next.ServeHTTP(w, r)
	})
}

// CreateSocket : Create a new server socket
func (socket *SocketConnection) CreateSocket() {
	socket.Server = socketio.NewServer(nil)
	socket.Controllers = []controller{}
	socket.Controllers = append(socket.Controllers, &GameRoomController{})
	socket.Controllers = append(socket.Controllers, &ChatController{})

	for _, controller := range socket.Controllers {
		controller.Init(socket.Server)
	}
}
