package services

import (
	"testing"

	"message-app/conf"
	mockServices "message-app/mocks/services"
	"message-app/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AuthenticationServiceTestSuite struct {
	suite.Suite
	mockCtrl        *gomock.Controller
	mockJwtService  *mockServices.MockJWTService
	mockConfService *mockServices.MockGbeConfigService
	authService     AuthenticationService
}

func (s *AuthenticationServiceTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())

	s.mockJwtService = mockServices.NewMockJWTService(s.mockCtrl)
	s.mockConfService = mockServices.NewMockGbeConfigService(s.mockCtrl)
	s.authService = NewAuthenticationService(s.mockJwtService, s.mockConfService)
}

func (s *AuthenticationServiceTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestAuthenticationService(t *testing.T) {
	suite.Run(t, new(AuthenticationServiceTestSuite))
}

func (s *AuthenticationServiceTestSuite) TestAuthenticationService() {
	s.Run("Should generate JWT token", func() {
		s.mockConfService.EXPECT().GetConfig().Return(conf.GetTestConf()).Times(1)
		s.mockJwtService.EXPECT().CreateLoginToken("Test").Return("Token", nil).Times(1)

		token, err := s.authService.CreateAuthToken(&models.AuthenticationReq{
			Password: "Test",
		})
		s.NoError(err)
		s.Equal("Token", token)
	})

	s.Run("Should not generate JWT token", func() {
		s.mockConfService.EXPECT().GetConfig().Return(conf.GetTestConf()).Times(1)

		_, err := s.authService.CreateAuthToken(&models.AuthenticationReq{
			Password: "Test-1",
		})
		s.Error(err)
	})

}
