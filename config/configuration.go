package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/TheThingsNetwork/ttn/core/types"
	"gopkg.in/yaml.v2"
)

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
}

func (e *EUI) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var eui string

	err := unmarshal(&eui)

	fmt.Println(err)
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

func NewConfig() *Config {
	var path string
	flag.StringVar(&path, "config", "./default-config.yaml", "path to config file")

	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := yaml.NewDecoder(file)

	var config Config
	decoder.Decode(&config)

	return &config
}
