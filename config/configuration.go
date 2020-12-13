package config

import (
	"flag"
	"os"

	"github.com/TheThingsNetwork/ttn/core/types"
	"gopkg.in/yaml.v2"
)

// Type-Alias for ttnsdk.AppEUI
type EUI types.AppEUI

type Config struct {
	Mqtt struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"mqtt"`
	TTN struct {
		AccessKey string `yaml:"access_key"`
		AppID     string `yaml:"app_id"`
		AppEUI    EUI    `yaml:"app_eui"`
	} `yaml:"ttn"`
	Ditto struct {
		Host           string `yaml:"host"`
		Namespace      string `yaml:"namespace"`
		ConnectionName string `yaml:"connection_name"`
	} `yaml:"ditto"`
	MetricsEndpoint string `yaml:"metrics_endpoint"`
}

// Custom Transform-Function for app_eui to get directy a types.AppUI
func (e *EUI) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var eui string

	err := unmarshal(&eui)

	if err != nil {
		return err
	}

	parsed, err := types.ParseAppEUI(eui)

	if err != nil {
		return err
	}

	*e = EUI(parsed)

	return nil
}

func NewConfig() (*Config, error) {
	var path string
	flag.StringVar(&path, "config", "./config.yaml", "path to config file")

	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := yaml.NewDecoder(file)

	var config Config
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
