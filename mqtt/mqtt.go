package mqtt

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aaronschweig/twinevent-ttn/ttn"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	MQTT_BROKER   = os.Getenv("MQTT_BROKER")
	MQTT_USER     = os.Getenv("MQTT_USER")
	MQTT_PASSWORD = os.Getenv("MQTT_PASSWORD")
	reg           = regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
)

func Start() MQTT.Client {

	if len(MQTT_BROKER) == 0 {
		MQTT_BROKER = "mq.jreichwald.de:1883"
	}

	if len(MQTT_USER) == 0 {
		MQTT_USER = "twinevent"
	}

	if len(MQTT_PASSWORD) == 0 {
		MQTT_PASSWORD = "twinevent"
	}

	opt := MQTT.NewClientOptions().AddBroker(MQTT_BROKER).SetUsername(MQTT_USER).SetPassword(MQTT_PASSWORD)

	client := MQTT.NewClient(opt)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func RegistrationHandler(c MQTT.Client, m MQTT.Message) {
	topic := m.Topic()

	match := reg.FindAllString(topic, -1)

	if len(match) != 1 {
		log.Print("Cloud not extract MAC-Adress from topic")
		return
	}

	mac := match[0]

	mac = strings.ReplaceAll(mac, ":", "-")

	log.Printf("Extracted MAC-Adress %s", mac)

	device, err := ttn.Get(mac)

	if err != nil {
		log.Printf("%s\n", err)
		log.Printf("Creating Device for %s\n", mac)

		ttn.RegisterDevice(mac, "")

		device, _ = ttn.Get(mac)
	}

	// TODO: Respond with config.json
	log.Printf("%#v \n", device)
}
