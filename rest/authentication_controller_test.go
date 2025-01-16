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

type AuthControllerTestSuite struct {
	suite.Suite
	mockCtrl        *gomock.Controller
	mockAuthService *mock_service.MockAuthenticationService
	rr              *httptest.ResponseRecorder
	router          *mux.Router
	authController  AuthenticationController
}

func (s *AuthControllerTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockAuthService = mock_service.NewMockAuthenticationService(s.mockCtrl)
	s.authController = NewAuthenticationController(s.mockAuthService)
	s.router = mux.NewRouter()
	s.rr = httptest.NewRecorder()
}

func TestAuthController(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (s *AuthControllerTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *AuthControllerTestSuite) TestCreateToken() {
	s.Run("Should auth token", func() {
		authReq := models.AuthenticationReq{
			Password: "test",
		}
		// Create JSON body
		body, err := json.Marshal(authReq)
		if err != nil {
			s.Fail("error in marshalling body")
		}
		s.mockAuthService.EXPECT().CreateAuthToken(&authReq).Return("Token", nil).Times(1)

		s.router.HandleFunc("/api/v1/auth", s.authController.CreateToken()).Methods("POST")

		request, _ := http.NewRequest(http.MethodPost, "/api/v1/auth", bytes.NewReader(body))
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
