package controllers

import (
	"net/http"
	"strconv"
	"task-manager/models"
	"time"
	"github.com/gin-gonic/gin"
	"task-manager/database" // Import the new database package
)


func GetCompletedTasks(c *gin.Context) {
	var tasks []models.Task
	database.DB.Where("status = ?", "completed").Find(&tasks) // Предположительно у вас есть статус для задач

	c.HTML(http.StatusOK, "task_view.html", gin.H{
		"tasks": tasks,
		// Передайте дополнительные данные
	})
}

func GetTomorrowTasks(c *gin.Context) {
	userIDStr := c.Query("user_id")        // Получить user_id из параметров запроса
	userID, err := strconv.Atoi(userIDStr) // Преобразование в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	var tasks []models.Task
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02") // Получаем завтрашнюю дату
	database.DB.Where("date(created_at) = ? and user_id = ?", tomorrow, userID).Find(&tasks)

	c.HTML(http.StatusOK, "task_view.html", gin.H{
		"tasks": tasks,
		// Передайте дополнительные данные
	})
}

func GetAllTasks(c *gin.Context) {
	userIDStr := c.Query("user_id")        // Получить user_id из параметров запроса
	userID, err := strconv.Atoi(userIDStr) // Преобразование в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	var tasks []models.Task
	//database.DB.Find(&tasks) // Получаем все задачи
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	c.HTML(http.StatusOK, "task_view.html", gin.H{
		"tasks": tasks,
		// Вы можете передать дополнительные данные, такие как имя пользователя и статус
	})
}

func GetTodayTasks(c *gin.Context) {
	userIDStr := c.Query("user_id")        // Получить user_id из параметров запроса
	userID, err := strconv.Atoi(userIDStr) // Преобразование в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	var tasks []models.Task
	today := time.Now().Format("2006-01-02")                                     // Форматируем текущую дату
	database.DB.Where("date(created_at) = ? and user_id = ?", today, userID).Find(&tasks) // В зависимости от вашей модели

	c.HTML(http.StatusOK, "task_view.html", gin.H{
		"tasks": tasks,
		// Передайте дополнительные данные
	})
}