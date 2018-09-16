package main

import (
	"fmt"
	"log"
	"time"

	"sandbox-influxdb/modules"
	"sandbox-influxdb/modules/weather"

	"github.com/joho/godotenv"
)

func init() {

}
func doEvery(d time.Duration, providers []weather.Provider) {
	for x := range time.Tick(d) {
		fmt.Printf("this time is %v\n", x)
		for _, p := range providers {
			p.Create()
		}
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	modules.Init()

	wea := weather.New()
	doEvery(3*time.Hour, wea)
}
