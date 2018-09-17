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

type Forecast struct {
	COD     string         `json:"cod"`
	Message float64        `json:"message"`
	Cnt     int            `json:"cnt"`
	List    []ListForecast `json:"list"`
	City    City           `json:"city"`
}

// ListForecast struct
type ListForecast struct {
	Dt      int             `json:"dt"`
	Main    MainList        `json:"main"`
	Weather []ListWeather   `json:"weather"`
	Clouds  ListClouds      `json:"clouds"`
	Wind    ListWind        `json:"wind"`
	Rain    ListRain        `json:"rain"`
	Sys     ForecastListSys `json:"sys"`
	DtTxt   string          `json:"dt_txt"`
}

type ForecastListSys struct {
	Pod string `json:"pod"`
}

const apiForecast = "https://api.openweathermap.org/data/2.5/forecast"

// Save func
func (forecast Forecast) Save() {
	fmt.Println("inserting forecast weather...")
	startTime := time.Now()
	noWorkers := 10

	var models = make(chan model, 100)
	go forecast.fetchData(models)
	workerPool(noWorkers, models)

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
}

func (forecast Forecast) fetchData(models chan model) {
	for _, val := range weatherCities {
		// var current Current
		url := apiForecast + "?id=" + strconv.Itoa(val.ID) + "&units=metric&appid=" + modules.WeatherAPIKey
		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &forecast)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		forecast.City.Name = val.Name

		var m model
		m = forecast
		models <- m
	}
	close(models)
}

func (forecast Forecast) insert() {
	c, bp := modules.NewConnection()

	for _, data := range forecast.List {
		icon := data.Weather[0].Icon[:2]
		iconValue, _ := strconv.Atoi(icon)

		tags := map[string]string{
			"api":      "forecast",
			"cityID":   strconv.Itoa(forecast.City.ID),
			"cityName": forecast.City.Name,
		}
		fields := map[string]interface{}{
			"status": iconValue,
		}
		dt := time.Unix(int64(data.Dt), 0)
		pt, err := client.NewPoint("weather", tags, fields, dt)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
