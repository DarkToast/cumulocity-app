package mqtt

import (
	"encoding/binary"
	"errors"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"tarent.de/schmidt/client-user/domain"
	"testing"
)

func TestConnectionErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		token         MQTT.Token
		errorExpected bool
	}{
		{"Correct connection", MockToken{ErrorValue: nil}, false},
		{"Failed connection", MockToken{ErrorValue: errors.New("some: error")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			makeMockControlData(tt.token)
			port := Port{Client: MockClient{}}
			err := port.Connect()

			if (err == nil) == tt.errorExpected {
				if tt.errorExpected {
					t.Error("Connect(): error expected was no one was given.")
				} else {
					t.Error("Connect(): error was not expected, but an error occured.")
				}
			}
		})
	}
}

func TestConnectionEstablished(t *testing.T) {
	t.Run("Connection already established", func(t *testing.T) {
		makeMockControlData(MockToken{ErrorValue: nil})

		port := Port{Client: MockClient{}}
		_ = port.Connect()
		_ = port.Connect()

		if mockControlData.MethodCounter["connect"] != 1 {
			t.Error("Connect(): An established connection should not reconnect")
		}
	})
}

func TestSubscriptionErrorHandling(t *testing.T) {
	makeMockControlData(MockToken{ErrorValue: nil})
	port := Port{Client: MockClient{}}
	_ = port.Connect()

	t.Run("Handle client error", func(t *testing.T) {
		makeMockControlData(MockToken{ErrorValue: errors.New("some: error")})
		_, err := port.Subscribe("/topic/main/#", 1)

		if err == nil {
			t.Error("Subscribe(): error expected was no one was given.")
		}
	})
}

func TestSubscription(t *testing.T) {
	type given struct {
		payload []byte
		topic   string
	}

	type expected struct {
		messageExpected bool
		temperature     domain.Temperature
		humidity        domain.Humidity
		deviceId        domain.DeviceId
	}

	tests := []struct {
		name     string
		given    given
		expected expected
	}{
		{"Correct data",
			given{valuesToByteArray(3250, 6420), "/d/42/th"},
			expected{true, 32.50, 64.20, 42},
		},
		{"Correct data with zeros",
			given{valuesToByteArray(3250, 0), "/d/42/th"},
			expected{true, 32.5, 0, 42},
		},
		{"Unparsable topic",
			given{valuesToByteArray(3250, 6420), "/d/foobar/th"},
			expected{false, 0, 0, 0},
		},
		{"Unparsable payload",
			given{[]byte{4, 6}, "/d/56/th"},
			expected{false, 0, 0, 0},
		},
	}

	makeMockControlData(MockToken{ErrorValue: nil})
	port := Port{Client: MockClient{}}
	_ = port.Connect()
	channel, err := port.Subscribe("/d/+/th", 0)
	if err != nil {
		t.Errorf("Subscribe(): error while subscribing to topic.")
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockControlData.TestMessageHandler(nil, MockMessage{tt.given.topic, tt.given.payload})

			if tt.expected.messageExpected {
				message := <-channel
				if message.DeviceId != tt.expected.deviceId ||
					message.Temperature != tt.expected.temperature ||
					message.Humidity != tt.expected.humidity {

					t.Errorf("Subscription(): Expected - Temperature: %.2f = %.2f; Humidity: %.2f = %.2f; DeviceId: %d = %d",
						message.Temperature, tt.expected.temperature,
						message.Humidity, tt.expected.humidity,
						message.DeviceId, tt.expected.deviceId)
				}
			} else {
				if len(channel) > 0 {
					t.Errorf("Subscription(): Expected no message on topic %s and payload %s", tt.given.topic, tt.given.payload)
				}
			}
		})
	}
}

func valuesToByteArray(temperature int, humidity int) []byte {
	ta := make([]byte, 2)
	ha := make([]byte, 2)
	binary.LittleEndian.PutUint16(ta, uint16(temperature))
	binary.LittleEndian.PutUint16(ha, uint16(humidity))

	return append(ta, ha...)
}
