package payment

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres() (*gorm.DB, error) {
	SQL_CONNECTION_URL := os.Getenv("SQL_CONNECTION_URL")
	db, err := gorm.Open(postgres.Open(SQL_CONNECTION_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	// db, err := gorm.Open(postgres.Open(SQL_CONNECTION_URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return db, nil
}
