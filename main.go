package main

import (
	"log"
	"sandbox-influxdb/modules"
	"sandbox-influxdb/modules/weather"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	modules.Init()

	wea := weather.New()

	c := cron.New()
	c.AddFunc("0 */3 * * * ", func() { weather.Save(wea) }) // Run every 3 hour
	c.Run()
}
