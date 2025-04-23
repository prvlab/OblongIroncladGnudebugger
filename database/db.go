package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"task-manager/config"
	"task-manager/models" // Import the models package
)

var DB *gorm.DB

func InitDB(appConfig config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		appConfig.Database.Host, appConfig.Database.User, appConfig.Database.Password,
		appConfig.Database.DbName, appConfig.Database.Port)
	
	fmt.Printf("Подключение к базе данных: %s:%d\n", appConfig.Database.Host, appConfig.Database.Port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	
	fmt.Println("Успешное подключение к базе данных")
	
	// Автомиграция создаст таблицы, если их нет
	err = DB.AutoMigrate(&models.Task{}, &models.Users{}, &models.PasswordReset{})
	if err != nil {
		return fmt.Errorf("ошибка автомиграции: %v", err)
	}
	
	return nil
}
