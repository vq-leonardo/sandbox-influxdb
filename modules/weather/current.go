package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sandbox-influxdb/modules"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const apiCurrent = "https://api.openweathermap.org/data/2.5/weather"

// Current struct
type Current struct {
	Coord      Coord         `json:"coord"`
	Weather    []ListWeather `json:"weather"`
	Base       string        `json:"base"`
	Main       MainList      `json:"main"`
	Visibility int           `json:"visibility"`
	Wind       ListWind      `json:"wind"`
	Clouds     ListClouds    `json:"clouds"`
	Rain       ListRain      `json:"rain"`
	Dt         int           `json:"dt"`
	Sys        CurrentSys    `json:"sys"`
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Cod        int           `json:"codxs"`
}

// CurrentSys struct
type CurrentSys struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

// Save func
func (current Current) Save() {
	fmt.Println("inserting current weather...")
	startTime := time.Now()
	noWorkers := 10

	var models = make(chan model, 100)
	go current.fetchData(models)
	workerPool(noWorkers, models)

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
}

func (current Current) fetchData(models chan model) {
	for _, val := range weatherCities {
		// var current Current
		url := apiCurrent + "?id=" + strconv.Itoa(val.ID) + "&units=metric&appid=" + modules.WeatherAPIKey
		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &current)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		current.Name = val.Name

		var m model
		m = current
		models <- m
	}
	close(models)
}

func (current Current) insert() {
	c, bp := modules.NewConnection()

	icon := current.Weather[0].Icon[:2]
	iconValue, _ := strconv.Atoi(icon)

	tags := map[string]string{
		"api":      "current",
		"cityID":   strconv.Itoa(current.ID),
		"cityName": current.Name,
	}
	fields := map[string]interface{}{
		"status": iconValue,
	}
	dt := time.Unix(int64(current.Dt), 0)
	pt, err := client.NewPoint("weather", tags, fields, dt)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
