package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"tarent.de/schmidt/client-user/application"
)

var port *Port

func CreatePort(configuration *application.Configuration) *Port {
	if port != nil {
		return port
	}

	clientOpts := MQTT.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", configuration.MQTT_HOSTNAME, configuration.MQTT_PORT)).
		SetClientID("gateway").
		SetCleanSession(true)
	mqttClient := MQTT.NewClient(clientOpts)
	mqttPort := &Port{Client: mqttClient}

	err := mqttPort.Connect()
	if err != nil {
		log.Fatalf("Could not connect to MQTT! Error was: %s", err.Error())
	}

	port = mqttPort
	return port
}
