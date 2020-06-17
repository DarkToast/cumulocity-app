package cumulocity

import (
	"log"
	"tarent.de/schmidt/client-user/domain"
	"time"
)

var deviceIdMap = map[domain.DeviceId]Id{
	domain.DeviceId(1): Id("9636292"),
}

func processMeasurement(channel <-chan domain.Measurement, cumulocityClient *Client) {
	device, _ := cumulocityClient.GetDevice("9636292")
	log.Print(device)

	for measurement := range channel {
		log.Printf("Got a new measurement with temp: %.2f and humidity %.2f for device %d\n", measurement.Temperature, measurement.Humidity, measurement.DeviceId)

		temperatureMetric := Temperature{value: float64(measurement.Temperature), unit: C}
		humidityMetric := Humidity{value: float64(measurement.Humidity)}
		cumulocityDeviceId := deviceIdMap[measurement.DeviceId]

		measurement := Measurement{
			Source:  cumulocityDeviceId,
			Time:    time.Now(),
			Type:    "DHT22",
			Metrics: []Metric{temperatureMetric, humidityMetric},
		}

		err := cumulocityClient.SendMeasurement(measurement)
		if err != nil {
			log.Printf("Error while processing measurement: %s", err.Error())
		}
	}
}
