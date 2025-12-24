package payment

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres() (*gorm.DB, error) {
	SQL_CONNECTION_URL := os.Getenv("SQL_CONNECTION_URL")
	db, err := gorm.Open(postgres.Open(SQL_CONNECTION_URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return db, nil
}
