package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	Id        uint64    `gorm:"autoIncrement;column:id" json:"id"`
	Content   string    `gorm:"column:content" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (u *Message) BeforeCreate(*gorm.DB) error {
	u.CreatedAt = time.Now().In(time.UTC)
	return nil
}
