package mqtt

import (
	"encoding/json"
	"regexp"
	"strings"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
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

type DeviceConfig struct {
	TTN struct {
		AppEUI       string `json:"app_eui"`
		DevID        string `json:"dev_id"`
		DevClass     string `json:"dev_class"`
		DevEUI       string `json:"dev_eui"`
		AppKey       string `json:"app_key"`
		LoraNodeDr   int8   `json:"lora_node_dr"`
		NorthPort    int16  `json:"north_port"`
		SouthPort    int16  `json:"south_port"`
		ChannelClear int8   `json:"channel_clear"`
	} `json:"ttn"`
}

func NewDeviceConfig() *DeviceConfig {
	dc := &DeviceConfig{}
	dc.TTN.DevClass = "A"
	dc.TTN.LoraNodeDr = 5
	dc.TTN.NorthPort = 2
	dc.TTN.SouthPort = 99
	dc.TTN.ChannelClear = -90
	return dc
}

func (dc *DeviceConfig) WithTTNDevice(device *ttnsdk.Device) *DeviceConfig {
	dc.TTN.AppEUI = device.AppEUI.String()
	dc.TTN.DevID = device.DevID
	dc.TTN.DevEUI = device.DevEUI.String()
	dc.TTN.AppKey = device.AppID

	return dc
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

	dc := NewDeviceConfig().WithTTNDevice(device)

	pl, err := json.Marshal(dc)
	if err != nil {
		ms.log.Error("Could not serialize json", "error", err)
	}
	if token := c.Publish(topic+"/response", 0b0, true, pl); token.Wait() && token.Error() != nil {
		ms.log.Error("Could not publish to", "topic", topic+"/response", "message", token.Error())
	}
}
