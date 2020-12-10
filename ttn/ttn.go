package ttn

import (
	"log"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	"github.com/TheThingsNetwork/go-utils/random"
	"github.com/TheThingsNetwork/ttn/core/types"
	"github.com/aaronschweig/twinevent-ttn/config"
)

type TTN struct {
	config        *config.Config
	deviceManager ttnsdk.DeviceManager
}

func NewTTN(cfg *config.Config) *TTN {
	return &TTN{
		cfg,
		nil,
	}
}

func (ttn *TTN) CreateConnection() ttnsdk.Client {

	ttnConfig := ttnsdk.NewCommunityConfig(ttn.config.TTN.AppID)

	client := ttnConfig.NewClient(ttn.config.TTN.AppID, ttn.config.TTN.AccessKey)

	devices, err := client.ManageDevices()

	ttn.deviceManager = devices

	if err != nil {
		log.Fatalf("%s: could not get device manager", ttn.config.TTN.AppID)
	}

	return client
}

func (ttn *TTN) Get(id string) (*ttnsdk.Device, error) {
	return ttn.deviceManager.Get(id)
}

func (ttn *TTN) RegisterDevice(mac string, description string) {

	device := new(ttnsdk.Device)
	device.AppID = ttn.config.TTN.AppID
	device.DevID = mac
	device.Description = description
	device.AppEUI = types.AppEUI(ttn.config.TTN.AppEUI)

	device.AppKey = new(types.AppKey)
	random.FillBytes(device.AppKey[:])

	var deviceEUI [8]byte
	random.FillBytes(deviceEUI[:])
	device.DevEUI = deviceEUI

	err := ttn.deviceManager.Set(device)

	if err != nil {
		log.Fatalf("%s: Could not create Device %#v", ttn.config.TTN.AppID, err)
	}
}
