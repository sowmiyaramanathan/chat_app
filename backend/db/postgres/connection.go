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
	Dbhost := os.Getenv("DB_HOST")
	Dbport := os.Getenv("DB_PORT")
	DbUser := os.Getenv("DB_USER")
	Dbname := os.Getenv("DB_NAME")
	Dbpassword := os.Getenv("DB_PASSWORD")
	Dbdriver := os.Getenv("DB_DRIVER")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", Dbhost, Dbport, DbUser, Dbname, Dbpassword)

	var logLevel logger.LogLevel

	if os.Getenv("ENV") == "PROD" {
		logLevel = logger.Warn
	} else {
		logLevel = logger.Info
	}

	Db, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		slog.Error("database connection failed", "driver", Dbdriver, "error", err)
		os.Exit(1)
	} else {
		slog.Info("database connected", "driver", Dbdriver)
	}
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
