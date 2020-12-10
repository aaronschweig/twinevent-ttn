package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronschweig/twinevent-ttn/config"
	"github.com/aaronschweig/twinevent-ttn/ditto"
	"github.com/aaronschweig/twinevent-ttn/mqtt"
	"github.com/aaronschweig/twinevent-ttn/ttn"
	"github.com/hashicorp/go-hclog"
)

func main() {
	log := hclog.Default()

	log.Info("Reading Configuration...")
	conf, err := config.NewConfig()

	if err != nil {
		log.Error("Could not read configuration.", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Initialize Services
	log.Info("Initialize Services...")
	ttnService := ttn.NewTTNService(conf)
	ms := mqtt.NewMqttService(ttnService)
	ds := ditto.NewDittoService(conf)

	// Start MQTT
	log.Info("Connect to MQTT-Broker...")
	client := ms.Start()
	defer client.Disconnect(250)

	// Start TTN
	log.Info("Connect to TTN...")
	ttnClient := ttnService.CreateConnection()
	defer ttnClient.Close()

	log.Info("Setup Subscription...")
	token := client.Subscribe("registration/+", 0b0, ms.RegistrationHandler)

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	log.Info("Ensure Ditto-TTN-Connection...")
	if err = ds.CreateTTNConnection(); err != nil {
		panic(err)
	}

	log.Info("Running...")
	<-c
}
