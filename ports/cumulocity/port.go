package cumulocity

import (
	"net/http"
	"tarent.de/schmidt/cumulocity-gateway/configuration"
	"tarent.de/schmidt/cumulocity-gateway/domain"
)

type Port struct {
	Measurements chan<- domain.Measurement
}

func CreatePort(configuration *configuration.Config, http *http.Client) *Port {
	httpClient := &HttpClient{config: configuration, httpClient: http}
	cumulocityClient := &Client{httpClient: httpClient}

	channel := make(chan domain.Measurement, 10)
	go processMeasurement(channel, cumulocityClient)
	return &Port{Measurements: channel}
}
