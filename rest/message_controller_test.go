package rest

import (
	"bytes"
	"encoding/json"
	"io"
	mock_service "message-app/mocks/services"
	"message-app/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type MessageControllerTestSuite struct {
	suite.Suite
	mockCtrl           *gomock.Controller
	mockMessageService *mock_service.MockMessageService
	mockJwtService     *mock_service.MockJWTService
	rr                 *httptest.ResponseRecorder
	router             *mux.Router
	messageController  MessageController
	middlware          MiddleWare
}

func (s *MessageControllerTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockMessageService = mock_service.NewMockMessageService(s.mockCtrl)
	s.mockJwtService = mock_service.NewMockJWTService(s.mockCtrl)
	s.messageController = NewMessageController(s.mockMessageService)
	s.router = mux.NewRouter()
	s.rr = httptest.NewRecorder()
	s.middlware = NewMiddleWare(s.mockJwtService)
}

func TestMessageController(t *testing.T) {
	suite.Run(t, new(MessageControllerTestSuite))
}

func (s *MessageControllerTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *MessageControllerTestSuite) TestCreateMessage() {
	s.Run("Should create message", func() {
		req := models.CreateMessageReq{
			Content: "test",
		}
		// Create JSON body
		body, err := json.Marshal(req)
		if err != nil {
			s.Fail("error in marshalling body")
		}
		s.mockJwtService.EXPECT().VerifyAuthToken("valid-token").Return("valid-uid", nil).Times(1)
		s.mockMessageService.EXPECT().CreateMessage(&req).Return(nil).Times(1)

		s.router.HandleFunc("/api/v1/message", s.messageController.HandlerWithMiddleware(s.middlware)).Methods("POST")

		request, _ := http.NewRequest(http.MethodPost, "/api/v1/message", bytes.NewReader(body))
		request.Header.Set("Authorization", "Bearer valid-token")
		request.Header.Set("Content-Type", "application/json")

		s.router.ServeHTTP(s.rr, request)

		s.Equal(200, s.rr.Code, "Status code is not 200")
		bodyBytes, err := io.ReadAll(s.rr.Body)
		if err != nil {
			s.Fail("error in reading body")
		}

		var standardResponse StandardResponse
		err = json.Unmarshal(bodyBytes, &standardResponse)
		if err != nil || !standardResponse.Result || standardResponse.Code != models.SUCCESS {
			s.Fail("fail")
		}

	})
}
