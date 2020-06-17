package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tarent.de/schmidt/client-user/application"
	"tarent.de/schmidt/client-user/ports/cumulocity"
	"tarent.de/schmidt/client-user/ports/mqtt"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// MQTT
	mqttPort := mqtt.CreatePort(application.ProductionConfig)
	mqttChannel, err := mqttPort.Subscribe(application.ProductionConfig.MQTT_TOPIC, 0)
	if err != nil {
		log.Fatalf("Could not subscribe to topic %s", application.ProductionConfig.MQTT_TOPIC)
	}

	// Cumulocity
	cumulocityPort := cumulocity.CreatePort(application.ProductionConfig, &http.Client{})

	// Application
	go application.Service(mqttChannel, cumulocityPort.Measurements)

	<-c
	mqttPort.Disconnect()
}
