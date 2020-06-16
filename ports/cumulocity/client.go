package cumulocity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"tarent.de/schmidt/client-user/application"
	"tarent.de/schmidt/client-user/domain"
	"time"
)

type Client struct {
	configuration *application.Configuration
	httpClient    *http.Client
}

func (client *Client) GetDevice(deviceId DeviceId) (*Device, error) {
	object, err := client.get(fmt.Sprintf("/inventory/managedObjects/%s", string(deviceId)))
	if err != nil {
		log.Printf("Error receving Device with id: %s", string(deviceId))
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, object["lastUpdated"].(string))
	if err != nil {
		log.Printf("Warning: DateTime could not be parsed. Value was: %s", object["lastUpdated"].(string))
		return nil, err
	}

	return &Device{
		Id:             DeviceId(object["id"].(string)),
		Name:           Name(object["name"].(string)),
		Owner:          Owner(object["owner"].(string)),
		Created:        t,
		ChildDeviceIds: []DeviceId{},
		ParentDeviceId: nil,
	}, nil
}

func (client *Client) FindByDomainDeviceId(id domain.DeviceId) (*Device, error) {
	return nil, nil
}

func (client *Client) SendMeasurement(measurement Measurement) error {
	_, err := client.post("/measurement/measurements", measurement.Json())
	if err != nil {
		log.Printf("Error while seding a new measurement for device: %s", string(measurement.Source))
		return err
	}
	return nil
}

func (client *Client) post(path string, body string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, client.configuration.COMULOCITY_URL+path, bytes.NewBufferString(body))
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, err
	}

	return client.request(req)
}

func (client *Client) get(path string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, client.configuration.COMULOCITY_URL+path, nil)
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, err
	}

	return client.request(req)
}

func (client *Client) request(req *http.Request) (map[string]interface{}, error) {
	req.SetBasicAuth(client.configuration.COMULOCITY_USERNAME, client.configuration.COMULOCITY_PASSWORD)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("An error occured: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading from stream: %s", err.Error())
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Printf("Error while parsing response JSON: %s", err.Error())
		return nil, err
	}

	return result, nil
}
