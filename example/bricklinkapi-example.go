package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	var envVar map[string]string

	envVar, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	blConsumerKey := envVar["CONSUMER_KEY"]
	blConsumerSecret := envVar["CONSUMER_SECRET"]
	blTokenValue := envVar["TOKEN_VALUE"]
	blTokenSecret := envVar["TOKEN_SECRET"]

	fmt.Println(blConsumerKey, blConsumerSecret, blTokenValue, blTokenSecret)
}
