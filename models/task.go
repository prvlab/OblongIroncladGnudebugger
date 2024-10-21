package models

import (
	"time"
)

// Task определяет структуру задачи
type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      uint64 `json:"user_id"` // ID пользователя Telegram
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string // Добавьте это поле
}
