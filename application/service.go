package application

import (
	"log"
	"tarent.de/schmidt/cumulocity-gateway/domain"
)

func Service(mqttPort <-chan domain.Measurement, cumulocityPort chan<- domain.Measurement) {
	for message := range mqttPort {
		cumulocityPort <- message
	}

	log.Println("Application stopped. Closing target channels")
	close(cumulocityPort)
}
