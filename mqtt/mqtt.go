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

type MqttService struct {
	ttn *ttn.TTN
}

func NewMqttService(ttn *ttn.TTN) *MqttService {
	return &MqttService{ttn}
}

func (ms *MqttService) Start() MQTT.Client {

	if len(MQTT_BROKER) == 0 {
		MQTT_BROKER = "mq.jreichwald.de:1883"
	}

	if len(MQTT_USER) == 0 {
		MQTT_USER = "twinevent"
	}

	if len(MQTT_PASSWORD) == 0 {
		MQTT_PASSWORD = "twinevent"
	}

	opt := MQTT.NewClientOptions().
		AddBroker(MQTT_BROKER).
		SetUsername(MQTT_USER).
		SetPassword(MQTT_PASSWORD).
		SetAutoReconnect(true)

	client := MQTT.NewClient(opt)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func (ms *MqttService) RegistrationHandler(c MQTT.Client, m MQTT.Message) {
	topic := m.Topic()

	match := reg.FindAllString(topic, -1)

	if len(match) != 1 {
		log.Print("Cloud not extract MAC-Adress from topic")
		return
	}

	mac := match[0]

	mac = strings.ReplaceAll(mac, ":", "-")

	log.Printf("Extracted MAC-Adress %s", mac)

	device, err := ms.ttn.Get(mac)

	if err != nil {
		log.Printf("%s\n", err)
		log.Printf("Creating Device for %s\n", mac)

		ms.ttn.RegisterDevice(mac, "")

		device, _ = ms.ttn.Get(mac)
	}

	// TODO: Respond with config.json
	log.Printf("%#v \n", device)
}
