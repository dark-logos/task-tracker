package tasks

import (
	"net/http"

	"task-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//! \fn GetTasksHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler to retrieve a user's tasks.
//! \param s Task service instance.
//! \return Gin handler function.
func GetTasksHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tasks, err := s.GetTasks(userID.(int))
		if err != nil {
			s.logger.Error("Failed to get tasks", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, tasks)
	}
}

//! \fn CreateTaskHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler to create a new task.
//! \param s Task service instance.
//! \return Gin handler function.
func CreateTaskHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			s.logger.Warn("Invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(&task); err != nil {
			s.logger.Warn("Validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		task.UserID = userID.(int)
		taskID, err := s.CreateTask(&task)
		if err != nil {
			s.logger.Error("Failed to create task", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Task created", "task_id": taskID})
	}
}

//! \fn GetTaskHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler to retrieve a specific task.
//! \param s Task service instance.
//! \return Gin handler function.
func GetTaskHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("id")
		userID, _ := c.Get("user_id")
		task, err := s.GetTask(taskID, userID.(int))
		if err != nil {
			s.logger.Warn("Task not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusOK, task)
	}
}

//! \fn UpdateTaskHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler to update a task.
//! \param s Task service instance.
//! \return Gin handler function.
func UpdateTaskHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("id")
		userID, _ := c.Get("user_id")
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			s.logger.Warn("Invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(&task); err != nil {
			s.logger.Warn("Validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := s.UpdateTask(&task, taskID, userID.(int))
		if err != nil {
			s.logger.Warn("Task not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
	}
}

//! \fn DeleteTaskHandler(s *Service) gin.HandlerFunc
//! \brief Creates a Gin handler to delete a task.
//! \param s Task service instance.
//! \return Gin handler function.
func DeleteTaskHandler(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("id")
		userID, _ := c.Get("user_id")
		err := s.DeleteTask(taskID, userID.(int))
		if err != nil {
			s.logger.Warn("Task not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
	}
}