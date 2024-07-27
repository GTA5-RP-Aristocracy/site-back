package db

import (
	"database/sql"
	"fmt"
)

func ConnectDB(user, pass, host, database string) (*sql.DB, error) {
	// Connect to the database.
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, host, database))
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
