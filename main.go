package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronschweig/twinevent-ttn/config"
	"github.com/aaronschweig/twinevent-ttn/ditto"
	"github.com/aaronschweig/twinevent-ttn/mqtt"
	"github.com/aaronschweig/twinevent-ttn/ttn"
	"github.com/hashicorp/go-hclog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	log.Info("Initializing Services...")
	ttnService := ttn.NewTTNService(conf, log)
	ds := ditto.NewDittoService(conf)
	ms := mqtt.NewMqttService(ttnService, ds, log, conf)

	// Start MQTT
	log.Info("Connecting to MQTT-Broker...")
	client := ms.Start()
	defer client.Disconnect(250)

	// Start TTN
	log.Info("Connecting to TTN...")
	ttnClient := ttnService.CreateConnection()
	defer ttnClient.Close()

	log.Info("Setting up Subscription...")
	token := client.Subscribe("registration/+", 0b10, ms.RegistrationHandler)

	if token.Wait() && token.Error() != nil {
		log.Error("Could not subscribe", "error", token.Error())
		os.Exit(1)
	}

	log.Info("Ensuring Ditto-TTN-Connection...")
	if err = ds.CreateTTNConnection(); err != nil {
		panic(err)
	}

	if len(conf.MetricsEndpoint) > 0 {
		go func() {
			http.Handle(conf.MetricsEndpoint, promhttp.Handler())
			http.ListenAndServe(":2112", nil)
		}()
	}

	log.Info("Running...")
	<-c
	os.Exit(0)
}
