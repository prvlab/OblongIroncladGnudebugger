package controllers

import (
	"net/http"
  "fmt"
	//"task-manager/models"
  "net/smtp"
	"github.com/gin-gonic/gin"
  //"github.com/jordan-wright/email"
  //"encoding/json"
  //"os"
  "task-manager/config"
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
  if err := sendEmailWithoutTLS(email.Email, password); err != nil {  
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить письмо с паролем"})  
    return  
  }  

  // Возвращаем успешный ответ с данными пользователя  
  c.JSON(http.StatusOK, gin.H{"message": "Письмо с инструкцией по сбросу пароля отправлено на ваш email. Проверьте папку входящих."})


}


func sendEmailWithoutTLS(to, password string) error {
    
  config, err := config.LoadConfig("config/config.json")  
  if err != nil {  
      fmt.Println("Error loading config:", err)  
      return err 
  }
  //from := "info@prvlab.ru"
  //pass := "fcvycYrUceBGVTLRxrQd"
  //smtpHost := "smtp.mail.ru"
  //smtpPort := "25" // Порт для SMTP без TLS (обычно 587 или 25)

  msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\nВаш пароль: %s", config.Emails.Email, to, config.Emails.Subject, password)
  auth := smtp.PlainAuth("", config.Emails.Email, config.Emails.Password, config.Emails.SmtpServer)

  // Отправка email без TLS - ОПАСНО!
  err = smtp.SendMail(config.Emails.SmtpServer+":"+ config.Emails.SmtpPort, auth, config.Emails.Email, []string{to}, []byte(msg))
  if err != nil {
    return fmt.Errorf("не удалось отправить email: %w", err)
  }
  return nil
}
