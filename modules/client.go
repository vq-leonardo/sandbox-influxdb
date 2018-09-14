package modules

import (
	"log"
	"os"

	client "github.com/influxdata/influxdb/client/v2"
)

const (
	MyDB     = "mydb"
	username = "leo"
	password = "leo123"
)

var WeatherAPIKey string

func Init() {
	WeatherAPIKey = os.Getenv("WEATHERKEY")
}

// NewClient func
func NewClient() (client.Client, client.BatchPoints) {
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
		Database:  MyDB,
		Precision: "us",
	})
	if err != nil {
		log.Fatal(err)
	}

	return c, bp
}
