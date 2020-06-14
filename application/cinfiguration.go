package application

type Configuration struct {
	MQTT_HOSTNAME string
	MQTT_PORT     int
	MQTT_TOPIC    string

	COMULOCITY_URL      string
	COMULOCITY_USERNAME string
	COMULOCITY_PASSWORD string
}

var ProductionConfig = &Configuration{
	MQTT_HOSTNAME:       "192.168.0.11",
	MQTT_PORT:           1883,
	MQTT_TOPIC:          "/d/+/th",
	COMULOCITY_URL:      "https://tarent-gmbh.cumulocity.com",
	COMULOCITY_USERNAME: "c.schmidt@tarent.de",
	COMULOCITY_PASSWORD: "mrMjPh8zNm7VmQn",
}
