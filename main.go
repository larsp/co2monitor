package main

import (
	"log"
	"net/http"

	"github.com/larsp/co2monitor/meter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	device     = kingpin.Arg("device", "CO2 Meter device, such as /dev/hidraw2").Required().String()
	listenAddr = kingpin.Arg("listen-address", "The address to listen on for HTTP requests.").
			Default(":8080").String()
)

var (
	temperature = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_temperature_celsius",
		Help: "Current temperature in Celsius",
	})
	co2 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_co2_ppm",
		Help: "Current CO2 level (ppm)",
	})
)

func init() {
	prometheus.MustRegister(temperature)
	prometheus.MustRegister(co2)
}

func main() {
	kingpin.Parse()
	http.Handle("/metrics", promhttp.Handler())
	go measure()
	log.Printf("Serving metrics at '%v/metrics'", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func measure() {
	meter := new(meter.Meter)
	err := meter.Open(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
		return
	}

	for {
		result, err := meter.Read()
		if err != nil {
			log.Fatalf("Something went wrong: '%v'", err)
		}
		temperature.Set(result.Temperature)
		co2.Set(float64(result.Co2))
	}
}
