package rest

import (
	"fmt"
	"message-app/services"
)

func StartServer(container *services.Container) *HttpServer {
	wsController := NewWebSocketController(container.Producer, container.Consumer, container.MessageService)
	middleWare := NewMiddleWare(container.JWTService)
	authenticationController := NewAuthenticationController(container.AuthenticationService)
	messageController := NewMessageController(container.MessageService)
	httpServer := NewHttpServer(
		container.GbeConfigService.GetConfig().Rest.Addr,
		middleWare,
		wsController,
		authenticationController,
		messageController,
	)

	go httpServer.Start()

	fmt.Println("rest server ok")
	return httpServer
}
