package main

import (
	"fmt"
	"log"
	"os"

	"encoding/json"

	"github.com/enrichman/gomuni"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	regionFolder := os.Getenv("REGION_FOLDER")
	cityFolder := os.Getenv("CITY_FOLDER")
	townFolder := os.Getenv("TOWN_FOLDER")

	country := gomuni.Load(regionFolder, cityFolder, townFolder)

	b, _ := json.Marshal(country)
	fmt.Println(string(b))
}
