package controllers

import (
	"fmt"
	"net/http"
  "time"
	//"task-manager/models"
	"net/smtp"

	"github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "github.com/google/uuid"
	//"github.com/jordan-wright/email"
	//"encoding/json"
	//"os"
	"task-manager/config"
	"task-manager/models"
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
func OkRegistrationPage(c *gin.Context) {
  c.HTML(http.StatusBadRequest, "ok.html", nil)
}

func ResetPasswordPage(c *gin.Context) {
  c.HTML(http.StatusBadRequest, "reset_password.html", nil)
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
		c.HTML(http.StatusBadRequest, "registration.html", gin.H{"message": "Пароли не совпадают"})
		return
	}
  
  // создаем новую задачу
  adduser := models.Users{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: user.Password}

  db.Create(&adduser)
	// Для примера отображаем данные на экране
	c.JSON(http.StatusOK, gin.H{"message": "Регистрация успешна!", "user": user})
}

func HandleResetPassword(c *gin.Context){
  var email Emails 
      
  if err := c.ShouldBindJSON(&email); err != nil{  
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})  
      return  
    }
// Проверяем email в базе данных. Если есть, то отправляем пароль на электронную почту. Если нет, то отправляем сообщение об ошибке.

  var users []models.Users
  //db.Where("email = ?", email.Email).Find(&users)
  if err := db.Where("email = ?", email.Email).First(&users).Error; err != nil {
    if err == gorm.ErrRecordNotFound {
      c.JSON(http.StatusBadRequest, gin.H{"message": "Пользователь с таким email не найден"})
      return
    }
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка базы данных"})
    return
  }

  // Генерируем уникальный токен
  token := uuid.New().String()
  // Создаем запись о сбросе пароля
 
  passwordReset := models.PasswordReset{
    UserID:    users[0].ID,
    Token:     token,
    ExpiresAt: time.Now().Add(time.Hour * 24), // Токен действителен 24 часа
  }
  if err := db.Create(&passwordReset).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка базы данных"})
    return
  }
  
  // Отправляем email со ссылкой для сброса пароля
  resetLink := fmt.Sprintf("https://mcic.events/reset_password?token=%s", token)
  if err := sendEmailWithoutTLS(email.Email, resetLink); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось отправить письмо"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "Письмо с инструкцией по сбросу пароля отправлено на ваш email. Проверьте папку входящих."})

}

func sendEmailWithoutTLS(to, resetLink string) error {
    
  config, err := config.LoadConfig("config/config.json")  
  if err != nil {  
      fmt.Println("Error loading config:", err)  
      return err 
  }

  msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\nСсылка для востановления пароля: %s", config.Emails.Email, to, config.Emails.Subject, resetLink)
  auth := smtp.PlainAuth("", config.Emails.Email, config.Emails.Password, config.Emails.SmtpServer)

  // Отправка email без TLS - ОПАСНО!
  err = smtp.SendMail(config.Emails.SmtpServer+":"+ config.Emails.SmtpPort, auth, config.Emails.Email, []string{to}, []byte(msg))
  if err != nil {
    return fmt.Errorf("не удалось отправить email: %w", err)
  }
  return nil
}
