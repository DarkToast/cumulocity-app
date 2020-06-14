package cumulocity

import (
	"fmt"
	"net/http"
	"tarent.de/schmidt/client-user/application"
	"tarent.de/schmidt/client-user/domain"
)

type Port struct {
	Measurements chan<- domain.Measurement
}

func CreatePort(configuration *application.Configuration, httpClient *http.Client) *Port {
	comulocityClient := &Client{configuration: configuration, httpClient: httpClient}

	channel := make(chan domain.Measurement, 10)
	go processMeasurement(channel, comulocityClient)
	return &Port{Measurements: channel}
}

func processMeasurement(channel <-chan domain.Measurement, comulocityClient *Client) {
	_, _ = comulocityClient.GetDevice(DeviceId("9636292"))
	for measurement := range channel {
		fmt.Printf("Got a new measurement with temp: %.2f and humidity %.2f\n", float64(measurement.Temperature), float64(measurement.Humidity))
	}
}
