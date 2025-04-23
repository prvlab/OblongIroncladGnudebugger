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
  "golang.org/x/crypto/bcrypt"
	//"github.com/jordan-wright/email"
	//"encoding/json"
	//"os"
	"task-manager/config"
	"task-manager/models"
  "task-manager/database"
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

type NewPasswordRequest struct {
  Token    string `json:"token" binding:"required"`
  Password string `json:"password" binding:"required"`
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

func NewPasswordPage(c *gin.Context){
  c.HTML(http.StatusBadRequest, "new_password.html", nil)
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

  // Хеширование нового пароля (используйте безопасный метод хеширования!)
  hashedPassword, err := hashPassword(user.Password) // Реализуйте функцию hashPassword
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка хеширования пароля"})
    return
  }
  user.Password = hashedPassword

  // создаем новую запись в таблице users
  adduser := models.Users{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Password: user.Password}

  if err := database.DB.Create(&adduser).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"message": "данная почта уже используется другим пользователем"})
      return
    }
  
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
  if err := database.DB.Where("email = ?", email.Email).First(&users).Error; err != nil {
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
  if err := database.DB.Create(&passwordReset).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка базы данных"})
    return
  }
  
  // Отправляем email со ссылкой для сброса пароля
  config, err := config.LoadConfig("config/config.json")  
  if err != nil {  
      fmt.Println("Error loading config:", err)  
      return 
  }
  resetLink := fmt.Sprintf("%s/new_password_page?token=%s", config.Emails.Link, token)
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
// NewPassword обрабатывает сброс пароля
func NewPassword(c *gin.Context) {
  var request NewPasswordRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
    return
  }

  var passwordReset models.PasswordReset
  if err := database.DB.Where("token = ?", request.Token).First(&passwordReset).Error; err != nil {
    if err == gorm.ErrRecordNotFound {
      c.JSON(http.StatusBadRequest, gin.H{"message": "Токен не найден или истек"})
      return
    }
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка базы данных"})
    return
  }

  // Проверка срока действия токена
  if passwordReset.ExpiresAt.Before(time.Now()) {
    c.JSON(http.StatusBadRequest, gin.H{"message": "Токен истек"})
    return
  }

  // Обновление пароля пользователя
  var user models.Users
  if err := database.DB.First(&user, passwordReset.UserID).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка базы данных"})
    return
  }

  // Хеширование нового пароля (используйте безопасный метод хеширования!)
  hashedPassword, err := hashPassword(request.Password) // Реализуйте функцию hashPassword
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка хеширования пароля"})
    return
  }
  user.Password = hashedPassword

  if err := database.DB.Save(&user).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления пароля"})
    return
  }

  // Удаление записи о сбросе пароля после успешного обновления
  if err := database.DB.Delete(&passwordReset).Error; err != nil {
    // Логирование ошибки, но не возвращаем ошибку клиенту, т.к. главное - обновление пароля
    fmt.Printf("Ошибка удаления записи о сбросе пароля: %v\n", err)
  }

  c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})
}

// hashPassword - функция хеширования пароля с использованием bcrypt
func hashPassword(password string) (string, error) {
  // Используем bcrypt.DefaultCost для баланса между безопасностью и производительностью.
  // Можно изменить на другое значение, но DefaultCost обычно является хорошим выбором.
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
  }
  return string(hashedPassword), nil
}


// comparePasswords - функция сравнения хешированного пароля с введенным
func comparePasswords(hashedPassword, password string) error {
  err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
  if err != nil {
    if err == bcrypt.ErrMismatchedHashAndPassword {
      return fmt.Errorf("неверный пароль")
    }
    return fmt.Errorf("ошибка сравнения паролей: %w", err)
  }
  return nil
}