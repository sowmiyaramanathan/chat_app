package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/controllers"
	e "backend/entities"
	"backend/models"
	"backend/routes"
	"backend/services"
	"backend/services/websocket"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	Dbhost := os.Getenv("DB_HOST")
	Dbport := os.Getenv("DB_PORT")
	DbUser := os.Getenv("DB_USER")
	Dbname := os.Getenv("DB_NAME")
	Dbpassword := os.Getenv("DB_PASSWORD")
	Dbdriver := os.Getenv("DB_DRIVER")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", Dbhost, Dbport, DbUser, Dbname, Dbpassword)
	Db, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot cannot to database %s, error occured - %s", Dbdriver, err)
	} else {
		log.Printf("We are connected to %s database", Dbdriver)
	}
	Db.Debug().AutoMigrate(&e.User{}, &e.Message{}, &e.Friends{})
	return Db
}

func Run() {
	Db := ConnectDatabase()
	m := models.New(Db)
	hub := websocket.NewHub()
	ws := websocket.New(hub)
	s := services.New(m, ws)
	c := controllers.New(s)
	r := routes.InitializeRoutes(c)

	fmt.Println("\nListening to port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
