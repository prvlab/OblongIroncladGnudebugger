package controllers

import (
	"net/http"
  "fmt"
	//"task-manager/models"
  "net/smtp"
	"github.com/gin-gonic/gin"
  //"encoding/json"
)

// User структура для представления пользователя
type User struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
// User структура для представления параметра  
type Emails struct {  
  Email string `json:"email"`
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

func ResetPasswordPage(c *gin.Context) {
  c.HTML(http.StatusBadRequest, "reset_password.html", nil)
}

func HandleResetPassword(c *gin.Context){

  
  var password ="45678903743"
  var email Emails 

  fmt.Sprintf("HandleResetPassword")
              
  if err := c.ShouldBindJSON(&email); err != nil{  
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})  
      return  
    }
  // Отправляем электронное письмо с паролем  
  if err := sendEmail(email.Email, password); err != nil {  
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить письмо с паролем"})  
    return  
  }  

  // Возвращаем успешный ответ с данными пользователя  
  c.JSON(http.StatusOK, gin.H{"message": "Письмо с инструкцией по сбросу пароля отправлено на ваш email. Проверьте папку входящих."})


}

// sendEmail отправляет email с паролем  
func sendEmail(to string, password string) error {  
  
  from := "r.pavlov@ano-mcms.ru"            // Ваш email  
  pass := "*Gector1003"                     // Ваш пароль от email  
  // Указываем SMTP сервер и порт  
  smtpHost := "mx.ano-mcms.ru"  
  smtpPort := "465" // Для TLS
  
  msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Ваш пароль\n\nВаш пароль: %s", from, to, password)  
  auth := smtp.PlainAuth("", from, pass, smtpHost)

  // Отправляем email  
  fmt.Println("Сервер %s\n", msg)
  return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg)) 
} 
