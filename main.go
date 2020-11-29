package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronschweig/twinevent-ttn/mqtt"
	"github.com/aaronschweig/twinevent-ttn/ttn"
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	client := mqtt.Start()
	defer client.Disconnect(250)

	ttnClient := ttn.CreateConnection()
	defer ttnClient.Close()

	token := client.Subscribe("aaron/+", byte(0), mqtt.RegistrationHandler)

	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	<-c
}
