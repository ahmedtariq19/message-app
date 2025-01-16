package postgres

import (
	"message-app/conf"
	"message-app/models"
	"message-app/repository"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MessageSuite struct {
	suite.Suite

	db   *gorm.DB
	repo repository.Store
}

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

func (s *MessageSuite) SetupSuite() {
	cfg := conf.GetTestConf()
	repo := SharedStore(cfg)
	s.Require().NotNil(store, "SharedStore should not return nil")

	s.db = store.GetDb()
	s.repo = store

	s.repo = repo
}

func (s *MessageSuite) SetupTest() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE messages RESTART IDENTITY CASCADE").Error)
}

func (s *MessageSuite) Test_CreateMessage() {
	msg := &models.Message{
		Content: "test-content",
	}

	result, err := s.repo.CreateMessage(msg)

	s.Require().NoError(err)
	s.Equal(msg.Content, result.Content)
}
