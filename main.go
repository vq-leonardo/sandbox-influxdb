package main

import (
	"log"

	"sandbox-influxdb/modules"
	"sandbox-influxdb/modules/weather"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	modules.Init()

	wea := weather.New()
	weather.Save(wea)
}
