package cumulocity

import (
	cumGen "github.com/tarent/gomulocity/generic"
	cumM "github.com/tarent/gomulocity/measurement"
	"net/http"
	"tarent.de/schmidt/cumulocity-gateway/configuration"
	"tarent.de/schmidt/cumulocity-gateway/domain"
)

type Port struct {
	Measurements chan<- domain.Measurement
}

func CreatePort(configuration *configuration.Config, http *http.Client) *Port {
	client := &cumGen.Client{
		HTTPClient: http,
		BaseURL:    configuration.Cumulocity.Url,
		Username:   configuration.Cumulocity.Username,
		Password:   configuration.Cumulocity.Password,
	}

	measurementApi := cumM.NewMeasurementApi(client)
	channel := make(chan domain.Measurement, 10)
	go processMeasurement(channel, measurementApi)
	return &Port{Measurements: channel}
}
