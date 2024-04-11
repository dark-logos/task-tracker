package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

//! \fn Connect(dbURL string) (*sql.DB, error)
//! \brief Establishes a connection to the database.
//! \param dbURL Database connection URL.
//! \return Database connection and error (if any).
func Connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}