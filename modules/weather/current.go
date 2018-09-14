package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sandbox-influxdb/modules"
	"strconv"
	"sync"
	"time"
)

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

// func (cr *Current) retrieveDatapoints()

// CurrentSys struct
type CurrentSys struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

// City struct
type City struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Coord   Coord  `json:"coord"`
}

// Coord struct
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// ListWeather struct
type ListWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// MainList struct
type MainList struct {
	Temp      float64 `json:"temp"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"frnd_level"`
	Humidity  float64 `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

// ListWind struct
type ListWind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

// ListClouds struct
type ListClouds struct {
	All int `json:"all"`
}

// ListRain struct
type ListRain struct {
	Hour float64 `json:"3h"`
}

var weatherCities []City

var weatherCurrents = make(chan Current, 100)

// var testWorkers = make(chan string, 100)

// var results = make(chan string, 100)

func init() {
	citiesWeather, err := ioutil.ReadFile("./json_sources/cities.json")
	// get openWeather cities
	if err != nil {
		panic(err)
	}
	json.Unmarshal(citiesWeather, &weatherCities)
}

func retrieveFromWeather() {
	for _, val := range weatherCities {
		var weatherCurrent Current
		url := "https://api.openweathermap.org/data/2.5/weather"
		url = url + "?id=" + strconv.Itoa(val.ID) + "&units=metric&appid=" + modules.WeatherAPIKey
		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(body, &weatherCurrent)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		weatherCurrent.Name = val.Name

		weatherCurrents <- weatherCurrent
	}
	close(weatherCurrents)
}

func worker(wg *sync.WaitGroup) {
	for current := range weatherCurrents {
		current.insert()
	}
	wg.Done()
}

func workerPool(noWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
}

func Create() {
	fmt.Println("creating")
	startTime := time.Now()
	noWorkers := 10

	go retrieveFromWeather()
	workerPool(noWorkers)

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")

}
