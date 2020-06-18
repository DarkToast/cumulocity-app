package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Mqtt struct {
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
		Topic    string `yaml:"topic"`
		ClientId string `yaml:"clientId"`
	} `yaml:"mqtt"`

	Cumulocity struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"cumulocity"`
}

func NewConfig() *Config {
	return (&Config{}).Load()
}

func (c *Config) Load() *Config {
	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("The configfile './config.yml' does not exist or is not readable. Concrete error: %s", err.Error())
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("YAML file was in the wrong format. Concrete error: %s", err.Error())
	}
	return c
}
