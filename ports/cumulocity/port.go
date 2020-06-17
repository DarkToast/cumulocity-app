package cumulocity

import (
	"net/http"
	"tarent.de/schmidt/client-user/application"
	"tarent.de/schmidt/client-user/domain"
)

type Port struct {
	Measurements chan<- domain.Measurement
}

func CreatePort(configuration *application.Configuration, http *http.Client) *Port {
	httpClient := &HttpClient{configuration: configuration, httpClient: http}
	cumulocityClient := &Client{httpClient: httpClient}

	channel := make(chan domain.Measurement, 10)
	go processMeasurement(channel, cumulocityClient)
	return &Port{Measurements: channel}
}
