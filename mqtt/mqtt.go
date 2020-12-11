package mqtt

import (
	"log"
	"regexp"
	"strings"

	"github.com/aaronschweig/twinevent-ttn/config"
	"github.com/aaronschweig/twinevent-ttn/ttn"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	reg = regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
)

type MqttService struct {
	ttn    *ttn.TTNService
	config *config.Config
}

func NewMqttService(ttn *ttn.TTNService, cfg *config.Config) *MqttService {
	return &MqttService{ttn, cfg}
}

func (ms *MqttService) Start() MQTT.Client {

	opt := MQTT.NewClientOptions().
		AddBroker(ms.config.Mqtt.Host).
		SetUsername(ms.config.Mqtt.User).
		SetPassword(ms.config.Mqtt.Password).
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
