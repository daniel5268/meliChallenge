package infrastructure

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDataSourceName() string {
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		name     = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
		password = os.Getenv("DB_PASSWORD")
	)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, name,
	)
}

// NewGormPostgresClient returns a gorm instance for the PostgreSQL DB
func NewGormPostgresClient() *gorm.DB {
	dsn := getDataSourceName()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection was not correctly established")
	}

	return db
}
