// internal/database/db.go
package database

import (
	"fmt"
	"log"
	"os"
	"skedda-goclone/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

// NewDatabase initializes a PostgreSQL connection using GORM
func NewDatabase() (*Database, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// Migrate applies schema migrations for all models
func (db *Database) Migrate() {
	// Register all models for migration here
	if err := db.AutoMigrate(&models.Teacher{}, &models.Student{}, &models.Booking{}, &models.Subject{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
