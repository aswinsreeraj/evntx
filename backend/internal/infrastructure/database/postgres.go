package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Calcutta",
		host, user, password, dbName, port, sslMode)

	if sslMode == "verify-full" {
		sslRootCert := os.Getenv("DB_SSL_ROOT_CERT")
		if sslRootCert == "" {
			sslRootCert = "./global-bundle.pem"
		}
		dsn += fmt.Sprintf(" sslrootcert=%s", sslRootCert)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db, nil
}
