package mqtt

import (
	"encoding/binary"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"regexp"
	"strconv"
	"tarent.de/schmidt/cumulocity-gateway/domain"
)

type Port struct {
	Client        MQTT.Client
	connected     bool
	subscriptions []chan domain.Measurement
}

type Message struct {
	Payload []byte
	Topic   string
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
			log.Printf("Closing MQTT channel %d", i)
			close(channel)
		}
	}
}

func (port *Port) Subscribe(topic string, qos int) (<-chan domain.Measurement, error) {
	channel := make(chan domain.Measurement, 10)

	if token := port.Client.Subscribe(topic, byte(qos), makeChannelReceiver(channel)); token.Wait() && token.Error() != nil {
		close(channel)
		return nil, token.Error()
	} else {
		log.Println("Subscribed to /d/1/th")
		port.subscriptions = append(port.subscriptions, channel)
		return channel, nil
	}
}

func (port *Port) Connected() bool {
	return port.connected
}

func makeChannelReceiver(c chan<- domain.Measurement) func(client MQTT.Client, message MQTT.Message) {
	return func(client MQTT.Client, message MQTT.Message) {
		log.Printf("Debug: Received message on topic %s", message.Topic())
		match := topicPattern.FindStringSubmatch(message.Topic())
		if len(match) < 2 {
			log.Printf("Error while parsing the numeric id of the mqtt topic: %s. Ignoring it! \n", message.Topic())
			return
		}

		if len(message.Payload()) != 4 {
			log.Print("Error while parsing the message payload. Expected a two uint16 (4 bytes) value.")
			return
		}

		foundId, _ := strconv.Atoi(match[1])
		deviceId := domain.DeviceId(foundId)
		temp := domain.Temperature(byteToFloat(message.Payload()[:2]))
		humidity := domain.Humidity(byteToFloat(message.Payload()[2:]))

		c <- domain.Measurement{
			Temperature: temp,
			Humidity:    humidity,
			DeviceId:    deviceId,
		}
	}
}

var topicPattern = regexp.MustCompile("^/d/([\\d]+)/th$")

func byteToFloat(array []byte) float64 {
	return float64(int16(binary.LittleEndian.Uint16(array))) / 100
}
