package cumulocity

import (
	"fmt"
	"log"
	"tarent.de/schmidt/client-user/domain"
	"time"
)

type Client struct {
	httpClient *HttpClient
}

func (client *Client) GetDevice(id Id) (*Device, error) {
	object, err := client.httpClient.get(fmt.Sprintf("/inventory/managedObjects/%s", string(id)))
	if err != nil {
		log.Printf("Error receving Device with id: %s", string(id))
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, object["lastUpdated"].(string))
	if err != nil {
		log.Printf("Warning: DateTime could not be parsed. Value was: %s", object["lastUpdated"].(string))
		return nil, err
	}

	return &Device{
		Id:            Id(object["id"].(string)),
		Name:          Name(object["name"].(string)),
		Owner:         Owner(object["owner"].(string)),
		CreationTime:  t,
		ChildDevices:  []Id{},
		ParentDevices: []Id{},
	}, nil
}

func (client *Client) FindByDomainDeviceId(id domain.DeviceId) (*Device, error) {
	return nil, nil
}

func (client *Client) SendMeasurement(measurement Measurement) error {
	_, err := client.httpClient.post("/measurement/measurements", measurement.Json())
	if err != nil {
		log.Printf("Error while seding a new measurement for device: %s", string(measurement.Source))
		return err
	}
	return nil
}
