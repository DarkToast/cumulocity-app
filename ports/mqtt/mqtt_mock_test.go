package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

// ----

type MockControlData struct {
	TestToken          MQTT.Token
	TestMessageHandler MQTT.MessageHandler
	TestTopic          string
	MethodCounter      map[string]int
}

var mockControlData MockControlData

func makeMockControlData(token MQTT.Token) {
	mockControlData = MockControlData{
		TestToken:          token,
		TestMessageHandler: nil,
		TestTopic:          "",
		MethodCounter:      make(map[string]int),
	}
}

// ----

type MockClient struct{}

func (client MockClient) IsConnected() bool      { return true }
func (client MockClient) IsConnectionOpen() bool { return true }
func (client MockClient) Connect() MQTT.Token {
	mockControlData.MethodCounter["connect"]++
	return mockControlData.TestToken
}
func (client MockClient) Disconnect(quiesce uint) {
	mockControlData.MethodCounter["disconnect"]++
}

func (client MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	mockControlData.TestTopic = topic
	return mockControlData.TestToken
}

func (client MockClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	mockControlData.MethodCounter["subscribe"]++
	mockControlData.TestMessageHandler = callback
	mockControlData.TestTopic = topic
	return mockControlData.TestToken
}

func (client MockClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	mockControlData.TestMessageHandler = callback
	return mockControlData.TestToken
}

func (client MockClient) Unsubscribe(topics ...string) MQTT.Token             { return mockControlData.TestToken }
func (client MockClient) AddRoute(topic string, callback MQTT.MessageHandler) {}
func (client MockClient) OptionsReader() MQTT.ClientOptionsReader             { return MQTT.ClientOptionsReader{} }

// ----

type MockToken struct {
	ErrorValue error
}

func (token MockToken) Wait() bool                     { return true }
func (token MockToken) WaitTimeout(time.Duration) bool { return true }
func (token MockToken) Error() error                   { return token.ErrorValue }

// ----

type MockMessage struct {
	TopicValue   string
	PayloadValue []byte
}

func (message MockMessage) Duplicate() bool   { return false }
func (message MockMessage) Qos() byte         { return 0 }
func (message MockMessage) Retained() bool    { return false }
func (message MockMessage) Topic() string     { return message.TopicValue }
func (message MockMessage) MessageID() uint16 { return 666 }
func (message MockMessage) Payload() []byte   { return message.PayloadValue }
func (message MockMessage) Ack()              {}
