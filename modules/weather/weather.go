package weather

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

// Provider interface
type Provider interface {
	Save()
}

type model interface {
	insert()
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

func init() {
	citiesWeather, err := ioutil.ReadFile("./json_sources/cities.json")
	// get openWeather cities
	if err != nil {
		panic(err)
	}
	json.Unmarshal(citiesWeather, &weatherCities)
}

func worker(wg *sync.WaitGroup, models chan model) {
	for m := range models {
		m.insert()
	}
	wg.Done()
}

func workerPool(noWorkers int, models chan model) {
	var wg sync.WaitGroup
	for i := 0; i < noWorkers; i++ {
		wg.Add(1)
		go worker(&wg, models)
	}
	wg.Wait()
}

// New func
func New() []Provider {
	var current Current
	var forecast Forecast
	var providers = []Provider{
		current,
		forecast,
	}
	return providers
}

// Save func
func Save(providers []Provider) {
	for _, p := range providers {
		p.Save()
	}
}
