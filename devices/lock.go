package devices

import (
	"fmt"
	"log"
	"net/http"

	"maas360api/internal/constants"
)

type DeviceActionResponse struct {
	Maas360DeviceID string `json:"maas360DeviceId"`
	ActionStatus    int    `json:"actionStatus"`
	ActionID        int    `json:"actionID"`
	Description     string `json:"description"`
}

// LockDevice sends a request to lock a specific device in MaaS360.
func LockDevice(serviceURL string, billingID string, deviceID string, maasToken string) error {
	if serviceURL == "" || billingID == "" || deviceID == "" || maasToken == "" {
		return fmt.Errorf("serviceURL, billingID, deviceID, and maasToken must not be empty")
	}

	url := fmt.Sprintf("%s/device-apis/devices/1.0/lockDevice/%s?deviceId=%s", serviceURL, billingID, deviceID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeForm)
	req.Header.Set(constants.AcceptHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.AuthorizationHeader, fmt.Sprintf(constants.MaaSTokenPrefix, maasToken))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}
	var response DeviceActionResponse
	if response.ActionStatus != 0 {
		return fmt.Errorf("failed to lock device: %s", response.Description)
	} else {
		log.Printf("Device %s lock scheduled successfully", deviceID)
		return nil
	}
}
