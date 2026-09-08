package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	e "backend/entities"
	"backend/metrics"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var sqlDB *sql.DB

func ConnectDatabase() *gorm.DB {
	var logLevel logger.LogLevel
	sslMode := "disable"

	if os.Getenv("ENV") == "PROD" {
		logLevel = logger.Warn
		sslMode = "require"
	} else {
		logLevel = logger.Info
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDriver := os.Getenv("DB_DRIVER")
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", dbHost, dbPort, dbUser, dbName, sslMode, dbPassword)

	Db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		slog.Error("database connection failed", "driver", dbDriver, "error", err)
		os.Exit(1)
	}

	slog.Info("database connected", "driver", dbDriver)

	Db.AutoMigrate(&e.User{}, &e.Message{}, &e.Friends{})

	sqlDB, err = Db.DB()
	if err != nil {
		slog.Error("failed to get database connection pool", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxOpenConns(100)
	metrics.StartRuntimeSampler(sqlDB)

	return Db
}

func Ping() error {
	return sqlDB.Ping()
}
