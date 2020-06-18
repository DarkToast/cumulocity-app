package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"tarent.de/schmidt/cumulocity-gateway/configuration"
)

var port *Port

func CreatePort(configuration *configuration.Config) *Port {
	if port != nil {
		return port
	}

	clientOpts := MQTT.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", configuration.Mqtt.Hostname, configuration.Mqtt.Port)).
		SetClientID(configuration.Mqtt.ClientId).
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
