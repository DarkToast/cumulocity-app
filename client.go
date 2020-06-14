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

	mqttPort := mqtt.CreatePort(application.ProductionConfig)

	mqttChannel, err := mqttPort.Subscribe("/d/+/th", 0)
	if err != nil {
		log.Fatal("Could not subscribe to topic /d/+/th")
	}

	cumulocityChannel := cumulocity.CreatePort(application.ProductionConfig, &http.Client{})
	go application.Service(mqttChannel, cumulocityChannel.Measurements)

	<-c
	mqttPort.Disconnect()
}
