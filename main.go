package main

import (
	"regexp"
	"strings"

	"github.com/ssimunic/gosensors"
)

type SensorTemp struct {
	Temperature  string
	CriticalTemp string
}

func NewSensorTemp(temp string) SensorTemp {
	r, _ := regexp.Compile(`[0-9]+.[0-9](\\.[0-9]+)?`)
	var currentTemp = r.FindAllString(temp, 1)[0]
	var criticTemp = r.FindAllString(temp, 2)[1]
	return SensorTemp{Temperature: currentTemp, CriticalTemp: criticTemp}
}

func main() {
	sensors, err := gosensors.NewFromSystem()

	var temps []SensorTemp

	// Iterate over chips
	for chip := range sensors.Chips {
		// Iterate over entries
		for _, value := range sensors.Chips[chip] {
			// If CPU or GPU, print out
			if strings.Index(value, "crit") > 0 {
				var item = NewSensorTemp(value)
				temps = append(temps, item)
			}
		}
	}
	println(temps[0].Temperature)

	if err != nil {
		panic(err)
	}

	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, sensors.JSON())
	// })
	// e.Logger.Fatal(e.Start(":1323"))
}
