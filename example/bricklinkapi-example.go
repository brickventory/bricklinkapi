package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brickventory/bricklinkapi"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env variables for accessing the BL Api
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// we need a new Bricklink
	bl := bricklinkapi.New(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("TOKEN_VALUE"),
		os.Getenv("TOKEN_SECRET"),
	)

	// try some simple query
	fmt.Println(bl.GetItem("part", "3001"))

	// try some more sufficticated query
	priceParams := make(map[string]string)
	priceParams["guide_type"] = "sold"
	priceParams["new_or_used"] = "U"
	priceParams["country_code"] = "DE"

	itemPrice, _ := bl.GetItemPrice("part", "3001", priceParams)
	fmt.Println(itemPrice)

	// some other querys
	fmt.Println(bl.GetItemImage("part", "3001", 0))
	fmt.Println(bl.GetColorList())
	fmt.Println(bl.GetColor(104))
	fmt.Println(bl.GetCategoryList())
	fmt.Println(bl.GetCategory(10))
}
