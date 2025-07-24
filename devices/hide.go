package devices

import (
	"encoding/json"
	"fmt"
	"log"
	"maas360api/internal/constants"
	"net/http"
)

func HideDevice(serviceURL string, billingID string, deviceID string, maasToken string) error {
	if serviceURL == "" || billingID == "" || deviceID == "" || maasToken == "" {
		return fmt.Errorf("serviceURL, billingID, deviceID, and maasToken must not be empty")
	}
	url := fmt.Sprintf("%s/device-apis/devices/1.0/hideDevice/%s?deviceId=%s", serviceURL, billingID, deviceID)
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
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}
	if response.ActionStatus != 0 {
		return fmt.Errorf("failed to hide device: %s", response.Description)
	}
	log.Printf("Device %s hidden successfully", deviceID)
	return nil
}
