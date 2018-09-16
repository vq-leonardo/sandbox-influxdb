package modules

import (
	"log"
	"os"

	"github.com/influxdata/influxdb/client/v2"
)

// var WeatherAPIKey string
var (
	dbName        string
	username      string
	password      string
	WeatherAPIKey string
)

// Init func
func Init() {
	WeatherAPIKey = os.Getenv("WEATHERKEY")
	dbName = os.Getenv("DBNAME")
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
}

// NewConnection func
func NewConnection() (client.Client, client.BatchPoints) {
	// Create new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: "us",
	})
	if err != nil {
		log.Fatal(err)
	}

	return c, bp
}
