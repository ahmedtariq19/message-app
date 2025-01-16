package rest

import (
	"encoding/json"
	"log"
	"message-app/models"
	"message-app/pkg/rabbitmq"
	"message-app/services"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketController interface {
	Handler() func(w http.ResponseWriter, r *http.Request)
	HandlerWithMiddleware(middleware MiddleWare) func(w http.ResponseWriter, r *http.Request)
}

type websocketController struct {
	producer       rabbitmq.Producer
	consumer       rabbitmq.Consumer
	messageService services.MessageService
}

func NewWebSocketController(producer rabbitmq.Producer, consumer rabbitmq.Consumer, messageService services.MessageService) WebSocketController {
	return &websocketController{
		producer:       producer,
		consumer:       consumer,
		messageService: messageService,
	}
}

func (p *websocketController) Handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		// Upgrade to WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to set WebSocket upgrade: %v", err)
			return
		}
		defer conn.Close()

		messageChan := make(chan string)

		// defer func() {
		// 	close(messageChan) // Ensure the channel is closed to avoid leaks
		// }()
		done := make(chan struct{})

		go p.readMessageFromSocket(conn, done)
		go p.consumer.Consume(models.WebsocketQueue, messageChan, done)
		go p.writeToSocket(conn, messageChan, done)

		// for {
		// 	_, message, err := conn.ReadMessage()
		// 	if err != nil {
		// 		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		// 			log.Println("WebSocket connection closed by client.")
		// 			return
		// 		}

		// 		if ne, ok := err.(net.Error); ok && ne.Timeout() {
		// 			log.Println("Read timeout, connection might be idle.")
		// 			continue
		// 		}

		// 		log.Println("Error reading message:", err)
		// 		return
		// 	}

		// 	if err := p.producer.Publish(models.WebsocketQueue, string(message)); err != nil {
		// 		log.Println("Error publishing message:", err)
		// 		continue
		// 	}
		// msg := <-messageChan
		// var m string
		// if err := json.Unmarshal([]byte(msg), &m); err != nil {
		// 	log.Println("Error unmarshalling message:", err)
		// 	continue
		// }

		// if err := p.messageService.CreateMessage(&models.CreateMessageReq{
		// 	Content: m,
		// }); err != nil {
		// 	log.Println("Error creating message:", err)
		// 	continue
		// }

		// log.Printf("Received and processed message: %s", msg)
		// if err := conn.WriteMessage(messageType, []byte(msg)); err != nil {
		// 	log.Println("Error writing message:", err)
		// }
		// }
		select {}
	}
}

func (p *websocketController) HandlerWithMiddleware(middleware MiddleWare) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.ValidateAuthToken(http.HandlerFunc(p.Handler())).
			ServeHTTP(w, r)
	}
}

func (p *websocketController) writeToSocket(ws *websocket.Conn, chanMessage chan string, done chan struct{}) {
	for {

		select {
		case <-done:
			log.Println("Stopping writeToSocket due to connection close.")
			return
		case message := <-chanMessage:

			var m string
			if err := json.Unmarshal([]byte(message), &m); err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			if err := p.messageService.CreateMessage(&models.CreateMessageReq{
				Content: m,
			}); err != nil {
				log.Println("Error creating message:", err)
				continue
			}
			if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Error writing message:", err)
			}
			log.Printf("Received and processed message: %s", message)
		default:
			continue
		}

	}
}

func (p *websocketController) readMessageFromSocket(ws *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket connection closed by client.")
				return
			}

			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				log.Println("Read timeout, connection might be idle.")
				continue
			}

			log.Println("Error reading message:", err)
			return
		}

		if err := p.producer.Publish(models.WebsocketQueue, string(message)); err != nil {
			log.Println("Error publishing message:", err)
			continue
		}

	}

}
