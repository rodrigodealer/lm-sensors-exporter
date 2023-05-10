package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ssimunic/gosensors"
)

type SensorTemp struct {
	Temperature  float64
	CriticalTemp string
}

func NewSensorTemp(temp string) SensorTemp {
	r, _ := regexp.Compile(`[0-9]+.[0-9](\\.[0-9]+)?`)
	var currentTemp = r.FindAllString(temp, 1)[0]
	currentTempFloat, _ := strconv.ParseFloat(currentTemp, 8)
	var criticTemp = r.FindAllString(temp, 2)[1]
	return SensorTemp{Temperature: currentTempFloat, CriticalTemp: criticTemp}
}

var tempGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "temperature_celsius",
		Help: "Temperature in celsius",
	},
)

func main() {
	sensors, err := gosensors.NewFromSystem()
	reg := prometheus.NewRegistry()

	reg.MustRegister(tempGauge)

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

	tempGauge.Set(temps[0].Temperature)

	if err != nil {
		panic(err)
	}
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	e := echo.New()
	e.GET("/metrics", echo.WrapHandler(promHandler))
	e.Logger.Fatal(e.Start(":9393"))
}
