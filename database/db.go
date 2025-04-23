
package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"task-manager/config"
)

var DB *gorm.DB

func InitDB(appConfig config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		appConfig.Database.Host, appConfig.Database.User, appConfig.Database.Password,
		appConfig.Database.DbName, appConfig.Database.Port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	return DB.AutoMigrate(&models.Task{}, &models.Users{}, &models.PasswordReset{})
}
