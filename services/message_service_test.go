package services

import (
	"errors"
	"testing"

	mockStore "message-app/mocks/store"
	"message-app/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MessageServiceTestSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	mockStore      *mockStore.MockStore
	messageService MessageService
}

func (s *MessageServiceTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())

	s.mockStore = mockStore.NewMockStore(s.mockCtrl)
	s.messageService = NewMessageService(s.mockStore)
}

func (s *MessageServiceTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestMessageService(t *testing.T) {
	suite.Run(t, new(MessageServiceTestSuite))
}

func (s *MessageServiceTestSuite) TestExchangeRate() {
	s.Run("Should store message", func() {
		msg := &models.Message{
			Content: "test message",
		}

		s.mockStore.EXPECT().CreateMessage(msg).Return(msg, nil).Times(1)

		err := s.messageService.CreateMessage(&models.CreateMessageReq{
			Content: "test message",
		})
		s.NoError(err)
	})

	s.Run("Should return error", func() {
		msg := &models.Message{
			Content: "test message",
		}

		s.mockStore.EXPECT().CreateMessage(msg).Return(msg, errors.New("failed")).Times(1)

		err := s.messageService.CreateMessage(&models.CreateMessageReq{
			Content: "test message",
		})
		s.Error(err)
	})
}
