package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronschweig/twinevent-ttn/config"
	"github.com/aaronschweig/twinevent-ttn/mqtt"
	"github.com/aaronschweig/twinevent-ttn/ttn"
)

func main() {

	conf := config.NewConfig()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	ttnService := ttn.NewTTN(conf)
	ms := mqtt.NewMqttService(ttnService)

	client := ms.Start()
	defer client.Disconnect(250)

	ttnClient := ttnService.CreateConnection()
	defer ttnClient.Close()

	token := client.Subscribe("registration/+", 0b0, ms.RegistrationHandler)

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	<-c
}
