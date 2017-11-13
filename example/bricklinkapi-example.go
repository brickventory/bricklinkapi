package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dirkolbrich/bricklinkapi"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env variables for accessing the BL Api
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bl := &bricklinkapi.Bricklink{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		Token:          os.Getenv("TOKEN_VALUE"),
		TokenSecret:    os.Getenv("TOKEN_SECRET"),
	}

	fmt.Println(bl.GetItem("part", "3004"))
}
