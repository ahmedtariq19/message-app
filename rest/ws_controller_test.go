package rest

// import (
// 	"message-app/mocks/pkg/rabbitmq"
// 	"message-app/mocks/services"
// 	"message-app/models"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/golang/mock/gomock"
// 	"github.com/gorilla/websocket"
// 	"github.com/stretchr/testify/suite"
// )

// type WebSocketControllerTestSuite struct {
// 	suite.Suite
// 	mockCtrl            *gomock.Controller
// 	mockProducer        *rabbitmq.MockProducer
// 	mockConsumer        *rabbitmq.MockConsumer
// 	mockMessageService  *services.MockMessageService
// 	webSocketController WebSocketController
// 	mockJwtService      *services.MockJWTService
// 	router              *http.ServeMux
// 	middlware           MiddleWare
// }

// func (s *WebSocketControllerTestSuite) SetupTest() {
// 	s.mockCtrl = gomock.NewController(s.T())
// 	s.mockProducer = rabbitmq.NewMockProducer(s.mockCtrl)
// 	s.mockConsumer = rabbitmq.NewMockConsumer(s.mockCtrl)
// 	s.mockMessageService = services.NewMockMessageService(s.mockCtrl)
// 	s.middlware = NewMiddleWare(s.mockJwtService)
// 	s.webSocketController = NewWebSocketController(s.mockProducer, s.mockConsumer, s.mockMessageService)
// 	s.router = http.NewServeMux()

// 	s.router.HandleFunc("/ws", s.webSocketController.Handler())
// }

// func (s *WebSocketControllerTestSuite) TearDownTest() {
// 	s.mockCtrl.Finish()
// }

// func TestWebSocketController(t *testing.T) {
// 	suite.Run(t, new(WebSocketControllerTestSuite))
// }

// func (s *WebSocketControllerTestSuite) TestWebSocketHandler() {
// 	s.mockProducer.EXPECT().Publish(models.WebsocketQueue, gomock.Any()).Return(nil).Times(1)
// 	s.mockConsumer.EXPECT().Consume(models.WebsocketQueue, gomock.Any(), gomock.Any()).DoAndReturn(func(queueName string, msgCh chan string, done chan struct{}) error {
// 		go func() {
// 			time.Sleep(100 * time.Millisecond)
// 			msgCh <- "Test Message"
// 		}()
// 		return nil
// 	}).Times(1)

// 	s.mockMessageService.EXPECT().CreateMessage(gomock.Any()).Return(nil).Times(1)

// 	server := httptest.NewServer(s.router)
// 	defer server.Close()

// 	wsURL := "ws" + server.URL[4:] + "/ws"
// 	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		s.Fail("Failed to connect to WebSocket: %v", err)
// 	}
// 	defer conn.Close()

// 	err = conn.WriteMessage(websocket.TextMessage, []byte("Client Message"))
// 	if err != nil {
// 		s.Fail("Failed to send message: %v", err)
// 	}

// 	_, message, err := conn.ReadMessage()
// 	if err != nil {
// 		s.Fail("Failed to read message: %v", err)
// 	}

// 	if string(message) != "Test Message" {
// 		s.Fail("Expected message 'Test Message', got '%s'", string(message))
// 	}
// }
