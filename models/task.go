package models

import (
	"time"
)

//! \struct Task
//! \brief Represents a task in the system.
type Task struct {
	ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    Title       string    `json:"title" validate:"required"`
    Description string    `json:"description"`
    Status      string    `json:"status" validate:"required,oneof=pending done in_progress"`
    Priority    int       `json:"priority" validate:"gte=1"`
    DueDate     time.Time `json:"due_date"`
    CreatedAt   time.Time `json:"created_at"`
}