package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tarent.de/schmidt/cumulocity-gateway/application"
	"tarent.de/schmidt/cumulocity-gateway/configuration"
	"tarent.de/schmidt/cumulocity-gateway/ports/cumulocity"
	"tarent.de/schmidt/cumulocity-gateway/ports/mqtt"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	config := configuration.NewConfig()

	// MQTT
	mqttPort := mqtt.CreatePort(config)
	mqttChannel, err := mqttPort.Subscribe(config.Mqtt.Topic, 0)
	if err != nil {
		log.Fatalf("Could not subscribe to topic %s", config.Mqtt.Topic)
	}

	// Cumulocity
	cumulocityPort := cumulocity.CreatePort(config, &http.Client{})

	// Application
	go application.Service(mqttChannel, cumulocityPort.Measurements)

	<-c
	mqttPort.Disconnect()
}
