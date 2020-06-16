package cumulocity

import (
	"fmt"
	"tarent.de/schmidt/client-user/domain"
)

func processMeasurement(channel <-chan domain.Measurement, comulocityClient *Client) {
	_, _ = comulocityClient.GetDevice(DeviceId("9636292"))
	for measurement := range channel {
		fmt.Printf("Got a new measurement with temp: %.2f and humidity %.2f for device %d\n", measurement.Temperature, measurement.Humidity, measurement.DeviceId)
	}
}
