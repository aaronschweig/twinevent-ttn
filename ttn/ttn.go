package ttn

import (
	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	"github.com/TheThingsNetwork/go-utils/random"
	"github.com/TheThingsNetwork/ttn/core/types"
	"github.com/aaronschweig/twinevent-ttn/config"
	"github.com/hashicorp/go-hclog"
)

type TTNService struct {
	config        *config.Config
	deviceManager ttnsdk.DeviceManager
	log           hclog.Logger
}

func NewTTNService(cfg *config.Config, log hclog.Logger) *TTNService {
	return &TTNService{
		cfg,
		nil,
		log,
	}
}

func (ttn *TTNService) CreateConnection() ttnsdk.Client {

	ttnConfig := ttnsdk.NewCommunityConfig(ttn.config.TTN.AppID)

	client := ttnConfig.NewClient(ttn.config.TTN.AppID, ttn.config.TTN.AccessKey)

	devices, err := client.ManageDevices()

	ttn.deviceManager = devices

	if err != nil {
		ttn.log.Error("Could not get device manager", "error", err)
	}

	return client
}

func (ttn *TTNService) Get(id string) (*ttnsdk.Device, error) {
	return ttn.deviceManager.Get(id)
}

func (ttn *TTNService) RegisterDevice(mac string, description string) {

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
		ttn.log.Error("Could not create Device.", "error", err)
	}
}
