package main

import (
	"backend/server"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error getting env, not coming through %v", err)
	} else {
		fmt.Print("We are getting env values")
	}

	server.Run()
}
