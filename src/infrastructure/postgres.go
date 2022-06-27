package infrastructure

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGormPostgresClient returns a gorm instance for the PostgreSQL DB
func NewGormPostgresClient() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection was not correctly established")
	}

	return db
}
