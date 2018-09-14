package main

import (
	"fmt"
	"log"
	"sandbox-influxdb/modules"
	"sandbox-influxdb/modules/weather"
	"time"

	"github.com/joho/godotenv"
)

func doEvery(d time.Duration, f func()) {
	for x := range time.Tick(d) {
		fmt.Printf("this time is %v\n", x)
		f()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	modules.Init()
	// doEvery(10*time.Second, weather.Create)
	weather.Create()
}
