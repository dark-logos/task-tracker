package models

import "time"

//! \struct User
//! \brief Represents a user in the system.
type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username" validate:"required"`
    PasswordHash string    `json:"password_hash" validate:"required"`
    Email        string    `json:"email" validate:"required,email"`
    CreatedAt    time.Time `json:"created_at"`
}