package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Port struct {
	Client        MQTT.Client
	connected     bool
	subscriptions []chan []byte
}

func (port *Port) Connect() error {
	if port.connected == true {
		log.Printf("MQTT Port already connected\n")
		return nil
	}

	if token := port.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Error: unable to connect to server\n")
		return token.Error()
	} else {
		fmt.Printf("Connected to server \n")
		port.connected = true
		return nil
	}
}

func (port *Port) Disconnect() {
	if port.connected == false {
		log.Println("MQTT Port not connected. Can not disconnect it.")
	} else {
		port.Client.Disconnect(500)
		log.Println("MQTT Port disconnected")
		for i, channel := range port.subscriptions {
			log.Printf("Closing channel %d", i)
			close(channel)
		}
	}
}

func (port *Port) Subscribe(topic string, qos int) (<-chan []byte, error) {
	channel := make(chan []byte, 10)

	if token := port.Client.Subscribe(topic, byte(qos), makeChannelReceiver(channel)); token.Wait() && token.Error() != nil {
		close(channel)
		return nil, token.Error()
	} else {
		log.Println("Subscribed to /d/1/th")
		port.subscriptions = append(port.subscriptions, channel)
		return channel, nil
	}
}

func makeChannelReceiver(c chan<- []byte) func(client MQTT.Client, message MQTT.Message) {
	return func(client MQTT.Client, message MQTT.Message) {
		payload := message.Payload()
		log.Printf("Debug: Received message on topic %s", message.Topic())
		c <- payload
	}
}
