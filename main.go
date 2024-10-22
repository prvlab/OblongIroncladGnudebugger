package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"task-manager/config"
	"task-manager/controllers"
	"task-manager/models"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	version = "1.0.0"
	build   = "10092024"
)
var db *gorm.DB
var bot *tgbotapi.BotAPI
var appConfig config.Config
var users = make(map[int]int) // для хранения пользовательского статуса
var mu sync.Mutex             // мьютекс для безопасного доступа к map в goroutines

 

// функция для обновления пользовательского статуса
func updateUserStatus(userID int, status int) {
	mu.Lock()
	defer mu.Unlock()
	users[userID] = status
}

// функция для получения статуса пользователя
func getUserStatus(userID int) int {
	mu.Lock()
	defer mu.Unlock()
	return users[userID]
}

func initDB(appConfig config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		appConfig.Database.Host, appConfig.Database.User, appConfig.Database.Password,
		appConfig.Database.DbName, appConfig.Database.Port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&models.Task{})
}
func main() {
	fmt.Println("version=", version)
	fmt.Println("build=", build)
	// Загрузка конфигурации
	var err error
	appConfig, err = config.LoadConfig("config/config.json") // Используйте функцию LoadConfig из пакета config
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Инициализация базы данных
	initDB(appConfig)
	controllers.SetDatabase(db) // Установка базы данных в контроллер
	// Инициализация бота
	bot, err = tgbotapi.NewBotAPI(appConfig.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Настройка Gin
	router := gin.Default()

	// путь к статическим файлам
	// <script src="/static/js/scripts.js"></script>
	// <link rel="stylesheet" href="/static/css/styles.css">
	router.Static("/static", "./static")
	// путь к шаблонам
	router.LoadHTMLGlob("templates/*")

	// Определение маршрутов
	router.GET("/", homePage)
	//	router.GET("/home", homePage)
	router.POST("/createtask", createTask)

	router.GET("/addTask", addTask)

	router.GET("/tasks", getTasks)
	router.GET("/tasks/all", controllers.GetAllTasks)
	router.GET("/tasks/today", controllers.GetTodayTasks)
	router.GET("/tasks/tomorrow", controllers.GetTomorrowTasks)
	router.GET("/tasks/completed", controllers.GetCompletedTasks)

	router.GET("/task_count", getTaskCount)

	router.GET("/registration", controllers.RegistrationPage)
	router.POST("/registration", controllers.HandleRegister)
	router.GET("/ok_registration", controllers.OkRegistrationPage)
	

	// Запуск веб-сервера
	serverAddress := appConfig.Server.ServerHost + ":" + appConfig.Server.ServerPort
	log.Println("Server is running on", serverAddress)
	err = http.ListenAndServe(serverAddress, router)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

	// Запуск бота
	go startBot()

	// Блокировка главного потока
	select {}
}

// homePage отображает главную страницу
func homePage(c *gin.Context) {
	userIDStr := c.Query("user_id") // Получить user_id из параметров запроса
	os := c.Query("os")             // Получаем информацию о ОС
	//	if isMobile == "true" {
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr) // Преобразование в int
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		// Использовать userID для получения информации о пользователе
		status := getUserStatus(userID)

		// Передать в HTML-ответ информацию о пользователе
		c.HTML(http.StatusOK, "index_vhod.html", gin.H{
			"userID":   userID,
			"status":   status,             // Проверьте свой статус пользователя
			"username": "Имя пользователя", // здесь вы можете настроить получение имени
			"os":       os,                 // Показать ОС
		})
	} else {
		c.HTML(http.StatusBadRequest, "index_vhod.html", nil)
	}
	//	} else {
	//		// Если не мобильное устройство, ничего не выводим
	//		c.HTML(http.StatusOK, "index.html", nil)
	//		//c.Status(http.StatusNoContent) // Вернуть 204 No Content
	//	}
}

func addTask(c *gin.Context) {
	c.HTML(http.StatusOK, "add_task.html", nil)
}

// createTask создает новую задачу
func createTask(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	userIDStr := c.PostForm("user_id") // userID как строка

	userIDStr = "1063764647" // для теста. убрать
	// Преобразование userID из строки в uint64
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	// создаем новую задачу
	task := models.Task{Title: title, Description: description, UserID: userID}

	db.Create(&task)

	c.Redirect(http.StatusFound, "/tasks")
}

// getTasks возвращает списки задач
func getTasks(c *gin.Context) {
	var tasks []models.Task
	db.Find(&tasks)

	c.HTML(http.StatusOK, "task_view.html", gin.H{
		"tasks": tasks,
	})
}

// startBot запускает Telegram-бота
func startBot() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // ignore non-Message Updates
			continue
		}
		userID := int64(update.Message.From.ID)
		updateUserStatus(int(userID), 1) // например, 1 для активного статуса
		// Обработка сообщений, команд и создание задач
		if update.Message.IsCommand() {
			// ... (добавьте свои команды здесь)
		} else {
			// создание задачи из сообщения (в случае нет команды)
			createTaskFromMessage(update.Message)
		}
	}
}

func createTaskFromMessage(message *tgbotapi.Message) {
	// Преобразуем userID в uint64
	userID := uint64(message.From.ID)

	task := models.Task{Title: message.Text, UserID: userID}
	db.Create(&task)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Task created: "+message.Text)
	bot.Send(msg)
}

func getTaskCount(c *gin.Context) {
	userIDStr := c.Query("user_id") // Получаем user_id из параметров запроса
	userID, err := strconv.Atoi(userIDStr)
	userID = 1063764647 // для теста. убрать
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	// Здесь вы должны получить количество всех задач из базы данных для данного userID
	allCount := getAllTaskCountForUser(userID) // Предположим, у вас есть такая функция

	c.JSON(http.StatusOK, gin.H{"count": allCount})
}

// Пример функции для получения количества задач
// принимает userID и возвращает количество задач для этого пользователя
func getAllTaskCountForUser(userID int) int64 {
	// Здесь добавьте логику для получения количества задач из вашей базы данных
	var count int64

	err := db.Model(&models.Task{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}
