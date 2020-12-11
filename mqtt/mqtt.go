package mqtt

import (
	"regexp"
	"strings"

	"github.com/aaronschweig/twinevent-ttn/config"
	"github.com/aaronschweig/twinevent-ttn/ditto"
	"github.com/aaronschweig/twinevent-ttn/ttn"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-hclog"
)

var (
	reg = regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
)

type MqttService struct {
	ttn    *ttn.TTNService
	ds     *ditto.DittoService
	config *config.Config
	log    hclog.Logger
}

func NewMqttService(ttn *ttn.TTNService, ds *ditto.DittoService, log hclog.Logger, cfg *config.Config) *MqttService {
	return &MqttService{ttn, ds, cfg, log}
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
		ms.log.Info("Cloud not extract MAC-Adress from topic")
		return
	}

	mac := match[0]

	mac = strings.ReplaceAll(mac, ":", "-")

	ms.log.Info("Extracted MAC-Adress", "mac", mac)

	device, err := ms.ttn.Get(mac)

	if err != nil {
		ms.log.Error("Device not found", err)
		ms.log.Info("Creating Device for", mac)

		ms.ttn.RegisterDevice(mac, "")

		device, _ = ms.ttn.Get(mac)
	}
	err = ms.ds.CreateDT(device)
	if err != nil {
		ms.log.Error("Error while creating DT", "error", err)
	}

	// TODO: Respond with config.json
	ms.log.Info("Device", device)
}
