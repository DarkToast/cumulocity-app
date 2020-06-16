package application

import (
	"encoding/binary"
	"tarent.de/schmidt/client-user/domain"
)

func byteToFloat(array []byte) float64 {
	return float64(binary.LittleEndian.Uint16(array)) / 100
}

func Service(mqttPort <-chan domain.Measurement, cumulocityPort chan<- domain.Measurement) {
	for message := range mqttPort {
		cumulocityPort <- message
	}
}
