package models
import (
	"time"
	"gorm.io/gorm" // Импортируйте вашу библиотеку GORM
)

type Users struct {
	gorm.Model
	//UserID      uint   `gorm:"primaryKey"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"password"`
}

// Структура для хранения данных о сбросе пароля
type PasswordReset struct {
	gorm.Model
	UserID    uint
	Token     string
	ExpiresAt time.Time
}