package ttn

import (
	"log"
	"os"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	"github.com/TheThingsNetwork/go-utils/random"
	"github.com/TheThingsNetwork/ttn/core/types"
)

var (
	TTN_ACCESS_KEY = os.Getenv("TTN_ACCESS_KEY")
	TTN_APP_ID     = os.Getenv("TTN_APP_ID")
	deviceManager  ttnsdk.DeviceManager
)

func CreateConnection() ttnsdk.Client {

	if len(TTN_ACCESS_KEY) == 0 {
		log.Fatal("No TTN_ACCESS_KEY present ...\n")
	}

	if len(TTN_APP_ID) == 0 {
		TTN_APP_ID = "dev_aaronschweig_ttn-test"
	}

	config := ttnsdk.NewCommunityConfig(TTN_APP_ID)

	client := config.NewClient(TTN_APP_ID, TTN_ACCESS_KEY)

	devices, err := client.ManageDevices()

	deviceManager = devices

	if err != nil {
		log.Fatalf("%s: could not get device manager", TTN_APP_ID)
	}

	return client
}

func Get(id string) (*ttnsdk.Device, error) {
	return deviceManager.Get(id)
}

func RegisterDevice(mac string, description string) {

	device := new(ttnsdk.Device)
	device.AppID = TTN_APP_ID
	device.DevID = mac
	device.Description = description
	device.AppEUI = types.AppEUI{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x03, 0x8E, 0xA1} // TODO: Replace

	device.AppKey = new(types.AppKey)
	random.FillBytes(device.AppKey[:])

	var deviceEUI [8]byte
	random.FillBytes(deviceEUI[:])
	device.DevEUI = deviceEUI

	err := deviceManager.Set(device)

	if err != nil {
		log.Fatalf("%s: Could not create Device %#v", TTN_APP_ID, err)
	}
}
