package config

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitDB() (*sql.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Retrieve database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Construct MySQL DSN (Data Source Name)
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ")/" + dbName + "?parseTime=true"

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
