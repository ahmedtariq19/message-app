package rest

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	addr                     string
	middleWare               MiddleWare
	wsController             WebSocketController
	authenticationController AuthenticationController
	messageController        MessageController
}

func NewHttpServer(
	addr string,
	middleware MiddleWare,
	websocketController WebSocketController,
	authenticationController AuthenticationController,
	messageController MessageController,
) *HttpServer {
	return &HttpServer{
		addr:                     addr,
		middleWare:               middleware,
		wsController:             websocketController,
		authenticationController: authenticationController,
		messageController:        messageController,
	}
}

func (server *HttpServer) Start() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	{
		api.HandleFunc("/auth", server.authenticationController.CreateToken()).Methods("POST")
	}

	ws := api.PathPrefix("/ws").Subrouter()
	{
		// Apply middleware to WebSocket handler
		// This makes sure token validation happens before WebSocket upgrade
		ws.HandleFunc("", server.wsController.HandlerWithMiddleware(server.middleWare))
	}
	message := api.PathPrefix("/message").Subrouter()
	{
		message.HandleFunc("", server.messageController.HandlerWithMiddleware(server.middleWare)).Methods("POST")

	}

	// Start server
	log.Printf("Starting server on %s", server.addr)
	err := http.ListenAndServe(server.addr, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
