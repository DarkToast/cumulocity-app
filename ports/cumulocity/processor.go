package cumulocity

import (
	m "github.com/tarent/gomulocity/measurement"
	"log"
	"tarent.de/schmidt/cumulocity-gateway/domain"
	"time"
)

var deviceIdMap = map[domain.DeviceId]m.Source{
	domain.DeviceId(1): {Id: "9636292"},
}

func processMeasurement(channel <-chan domain.Measurement, measurementApi m.MeasurementApi) {
	for measurement := range channel {
		log.Printf("Got a new measurement with temp: %.2f, humidity %.2f and pressure %.2f for device %d\n",
			measurement.Temperature, measurement.Humidity, measurement.AirPressure, measurement.DeviceId)

		source, ok := deviceIdMap[measurement.DeviceId]
		if !ok {
			log.Printf("Can not map device id %d to cumulocity source. Ignore measurement", measurement.DeviceId)
			break
		}

		t := m.ValueFragment{
			Value: float64(measurement.Temperature),
			Unit:  "C",
		}
		h := m.ValueFragment{
			Value: float64(measurement.Humidity),
			Unit:  "hPa",
		}
		p := m.ValueFragment{
			Value: float64(measurement.AirPressure),
			Unit:  "hPa",
		}
		now := time.Now()
		newMeasurement := &m.NewMeasurement{
			Time:            &now,
			MeasurementType: "",
			Source:          source,
			Temperature:     map[string]m.ValueFragment{"desk": t},
			Humidity:        map[string]m.ValueFragment{"desk": h},
			AirPressure:     map[string]m.ValueFragment{"desk": p},
		}

		_, err := measurementApi.Create(newMeasurement)
		if err != nil {
			log.Printf("Error while creating cumulocty measurement: %s", err.Error())
		}
	}
}
