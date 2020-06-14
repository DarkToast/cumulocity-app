package application

import (
	"encoding/binary"
	"tarent.de/schmidt/client-user/domain"
)

func byteToFloat(array []byte) float64 {
	return float64(binary.LittleEndian.Uint16(array)) / 100
}

func Service(mqttPort <-chan []byte, cumulocityPort chan<- domain.Measurement) {
	for message := range mqttPort {
		temp := domain.Temperature(byteToFloat(message[:2]))
		humidity := domain.Humidity(byteToFloat(message[2:]))
		measurement := domain.Measurement{
			Temperature: temp,
			Humidity:    humidity,
		}

		cumulocityPort <- measurement
	}
}
