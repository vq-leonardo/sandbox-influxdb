package weather

import (
	"log"
	"sandbox-influxdb/modules"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

func (weatherCurrent *Current) insert() {
	c, bp := modules.NewClient()

	icon := weatherCurrent.Weather[0].Icon[:2]
	iconValue, _ := strconv.Atoi(icon)

	tags := map[string]string{
		"api":      "current",
		"cityID":   strconv.Itoa(weatherCurrent.ID),
		"cityName": weatherCurrent.Name,
	}
	fields := map[string]interface{}{
		"status": iconValue,
	}
	dt := time.Unix(int64(weatherCurrent.Dt), 0)
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
