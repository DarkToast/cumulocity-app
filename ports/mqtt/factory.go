package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

var MQTT_HOSTNAME = "192.168.0.11"
var PORT = 1883
var port *Port

func CreatePort() *Port {
	if port != nil {
		return port
	}

	clientOpts := MQTT.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", MQTT_HOSTNAME, PORT)).
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
