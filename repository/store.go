package repository

import (
	"message-app/models"

	"gorm.io/gorm"
)

type Store interface {
	GetDb() *gorm.DB
	CreateMessage(msg *models.Message) (*models.Message, error)
}
