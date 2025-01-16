package services

import (
	"message-app/models"
	"message-app/repository"
)

type MessageService interface {
	CreateMessage(msg *models.CreateMessageReq) error
}

type messageService struct {
	store repository.Store
}

func NewMessageService(
	store repository.Store,
) MessageService {
	return &messageService{
		store: store,
	}
}

func (m *messageService) CreateMessage(req *models.CreateMessageReq) error {
	_, err := m.store.CreateMessage(&models.Message{
		Content: req.Content,
	})
	if err != nil {
		return err
	}

	return nil
}
