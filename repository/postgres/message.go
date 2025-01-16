package postgres

import "message-app/models"

func (s *Store) CreateMessage(msg *models.Message) (*models.Message, error) {
	if err := s.db.Model(&models.Message{}).Create(&msg).Error; err != nil {
		return nil, err
	}

	return msg, nil
}
