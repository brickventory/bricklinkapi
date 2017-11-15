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

	// we need a new Bricklink
	bl := bricklinkapi.New(os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("TOKEN_VALUE"),
		os.Getenv("TOKEN_SECRET"),
	)

	// try some simple query
	fmt.Println(bl.GetItem("part", "3004"))

	// try some more sufficticated query
	priceParams := make(map[string]string)
	priceParams["guide_type"] = "sold"
	priceParams["new_or_used"] = "U"
	priceParams["country_code"] = "DE"
	itemPrice, _ := bl.GetItemPrice("part", "3004", priceParams)
	fmt.Printf("got %v data sets\n", len(itemPrice))
	// print the first 100 bytes of the returned response body
	fmt.Println(itemPrice[:1000])
}
