package api

import socketio "github.com/googollee/go-socket.io"

// AuthenticationController : Controller for Google OAuth2
type AuthenticationController struct {
	socket *socketio.Server
}

// Init : Initialize the structure
func (controller *AuthenticationController) Init(socket *socketio.Server) {
	controller.socket = socket
	controller.socket.OnEvent("/", "msg", controller.GoogleAuth)
}

// GoogleAuth : Method to start the autentication
func (controller *AuthenticationController) GoogleAuth(socket *socketio.Server) {

}
