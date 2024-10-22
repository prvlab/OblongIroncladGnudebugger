package controllers

import (
	"net/http"
  //"fmt"
	//"task-manager/models"
	"github.com/gin-gonic/gin"
)

// User структура для представления пользователя
type User struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// Функция регистрации пользователя
func RegistrationPage(c *gin.Context) {
	c.HTML(http.StatusBadRequest, "registration.html", nil)
}

// handleRegister обрабатывает регистрацию пользователя
func HandleRegister(c *gin.Context) {
	var user User

	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Здесь можно добавить логику для проверки паролей и сохранения пользователя в базу данных
	if user.Password != user.ConfirmPassword {
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{"error": "Пароли не совпадают"})
		return
	}
  //fmt.Print("Регистрация успешна!", user)
	// Для примера отображаем данные на экране
	c.JSON(http.StatusOK, gin.H{"message": "Регистрация успешна!", "user": user})
}
func OkRegistrationPage(c *gin.Context) {
  c.HTML(http.StatusBadRequest, "ok.html", nil)
}
