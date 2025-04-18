package tasks

import (
	"database/sql"

	"task-tracker/internal/models"
	"go.uber.org/zap"
)

//! \struct Service
//! \brief Handles task-related business logic.
type Service struct {
	db     *sql.DB
	logger *zap.Logger
}

//! \fn NewService(db *sql.DB, logger *zap.Logger) *Service
//! \brief Initializes a new task service.
//! \param db Database connection.
//! \param logger Logger instance.
//! \return Pointer to initialized Service.
func NewService(db *sql.DB, logger *zap.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

//! \fn GetTasks(userID int) ([]models.Task, error)
//! \brief Retrieves tasks for a user.
//! \param userID ID of the user.
//! \return List of tasks and error (if any).
func (s *Service) GetTasks(userID int) ([]models.Task, error) {
	query := `SELECT id, user_id, title, description, status, priority, due_date, created_at 
              FROM tasks WHERE user_id = $1`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		s.logger.Error("Failed to fetch tasks", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, 
			&task.Status, &task.Priority, &task.DueDate, &task.CreatedAt); err != nil {
			s.logger.Error("Failed to scan task", zap.Error(err))
			continue
		}
		tasks = append(tasks, task)
	}

	s.logger.Info("Tasks retrieved", zap.Int("user_id", userID), zap.Int("count", len(tasks)))
	return tasks, nil
}

//! \fn CreateTask(task *models.Task) (int, error)
//! \brief Creates a new task in the database.
//! \param task Task data to create.
//! \return Task ID and error (if any).
func (s *Service) CreateTask(task *models.Task) (int, error) {
	query := `INSERT INTO tasks (user_id, title, description, status, priority, due_date) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var taskID int
	err := s.db.QueryRow(query, task.UserID, task.Title, task.Description, 
		task.Status, task.Priority, task.DueDate).Scan(&taskID)
	if err != nil {
		s.logger.Error("Failed to create task", zap.Error(err))
		return 0, err
	}

	s.logger.Info("Task created", zap.Int("task_id", taskID), zap.Int("user_id", task.UserID))
	return taskID, nil
}

//! \fn GetTask(taskID string, userID int) (*models.Task, error)
//! \brief Retrieves a specific task by ID for a user.
//! \param taskID ID of the task.
//! \param userID ID of the user.
//! \return Task and error (if any).
func (s *Service) GetTask(taskID string, userID int) (*models.Task, error) {
	var task models.Task
	query := `SELECT id, user_id, title, description, status, priority, due_date, created_at 
              FROM tasks WHERE id = $1 AND user_id = $2`
	err := s.db.QueryRow(query, taskID, userID).Scan(&task.ID, &task.UserID, &task.Title, 
		&task.Description, &task.Status, &task.Priority, &task.DueDate, &task.CreatedAt)
	if err != nil {
		s.logger.Warn("Task not found", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Task retrieved", zap.String("task_id", taskID), zap.Int("user_id", userID))
	return &task, nil
}

//! \fn UpdateTask(task *models.Task, taskID string, userID int) error
//! \brief Updates a task in the database.
//! \param task Updated task data.
//! \param taskID ID of the task.
//! \param userID ID of the user.
//! \return Error (if any).
func (s *Service) UpdateTask(task *models.Task, taskID string, userID int) error {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, due_date = $5 
              WHERE id = $6 AND user_id = $7`
	result, err := s.db.Exec(query, task.Title, task.Description, task.Status, task.Priority, 
		task.DueDate, taskID, userID)
	if err != nil {
		s.logger.Error("Failed to update task", zap.Error(err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Warn("Task not found", zap.String("task_id", taskID))
		return sql.ErrNoRows
	}

	s.logger.Info("Task updated", zap.String("task_id", taskID), zap.Int("user_id", userID))
	return nil
}

//! \fn DeleteTask(taskID string, userID int) error
//! \brief Deletes a task from the database.
//! \param taskID ID of the task.
//! \param userID ID of the user.
//! \return Error (if any).
func (s *Service) DeleteTask(taskID string, userID int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := s.db.Exec(query, taskID, userID)
	if err != nil {
		s.logger.Error("Failed to delete task", zap.Error(err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Warn("Task not found", zap.String("task_id", taskID))
		return sql.ErrNoRows
	}

	s.logger.Info("Task deleted", zap.String("task_id", taskID), zap.Int("user_id", userID))
	return nil
}