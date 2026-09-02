package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"backend/controllers"
	e "backend/entities"
	"backend/metrics"
	"backend/models"
	"backend/routes"
	"backend/services"
	"backend/services/websocket"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
		log.Fatalf("Cannot cannot to database %s, error occured - %s", Dbdriver, err)
	} else {
		log.Printf("We are connected to %s database", Dbdriver)
	}
	Db.AutoMigrate(&e.User{}, &e.Message{}, &e.Friends{})
	return Db
}

func Run() {
	Db := ConnectDatabase()
	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(100)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			stats := sqlDB.Stats()

			log.Printf(
				"DB_STATS open=%d in_use=%d idle=%d wait_count=%d wait_duration=%s",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration,
			)
		}
	}()
	metrics.StartRuntimeSampler(sqlDB)

	m := models.New(Db)
	hub := websocket.NewHub()
	ws := websocket.New(hub)
	s := services.New(m, ws)
	c := controllers.New(s)
	r := routes.InitializeRoutes(c)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Printf("METRICS_SNAPSHOT\n%s", metrics.Snapshot())
		}
	}()

	fmt.Println("\nListening to port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
