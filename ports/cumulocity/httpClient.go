package cumulocity

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"tarent.de/schmidt/client-user/application"
)

type HttpClient struct {
	configuration *application.Configuration
	httpClient    *http.Client
}

func (client *HttpClient) post(path string, body string) (map[string]interface{}, error) {
	url := client.configuration.COMULOCITY_URL + path

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, err
	}

	return client.request(req)
}

func (client *HttpClient) get(path string) (map[string]interface{}, error) {
	url := client.configuration.COMULOCITY_URL + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error: While creating a request: %s", err.Error())
		return nil, err
	}

	return client.request(req)
}

func (client *HttpClient) request(req *http.Request) (map[string]interface{}, error) {
	log.Printf("HTTP %s on URL %s", req.Method, req.URL)

	req.SetBasicAuth(client.configuration.COMULOCITY_USERNAME, client.configuration.COMULOCITY_PASSWORD)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("An error occured: %s", err.Error())
		return nil, err
	}
	log.Printf("Got status %d", resp.StatusCode)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading from stream: %s", err.Error())
		return nil, err
	}

	log.Printf("Debug: Response body was: %s", b)
	var result map[string]interface{}

	if len(b) > 0 {
		err = json.Unmarshal(b, &result)
		if err != nil {
			log.Printf("Error while parsing response JSON: %s", err.Error())
			return nil, err
		}
	}

	return result, nil
}
