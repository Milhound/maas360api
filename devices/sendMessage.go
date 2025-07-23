package devices_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	auth_api "maas360api/auth"
)

type sendMessageResponse struct {
	Maas360DeviceID string `json:"maas360DeviceID"`
	ActionStatus    int8   `json:"actionStatus"`
	Description     string `json:"description"`
}

type actionResponse struct {
	ActionResponse sendMessageResponse `json:"actionResponse"`
}

// SendMessage sends a message to a specific device in MaaS360.
// It requires a billing ID, device ID, an authentication token, and the message details.
func SendMessage(billingID string, deviceID string, messageTitle string, message string, maasToken string) error {
	if len(billingID) == 0 || len(deviceID) == 0 {
		return fmt.Errorf("billing ID and device ID cannot be empty")
	}
	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		return err
	}

	// Construct the message URL
	messageURL := fmt.Sprintf("%s/device-apis/devices/1.0/sendMessage/%s?deviceId=%s&messageTitle=%s&message=%s", serviceURL, billingID, deviceID, url.PathEscape(messageTitle), url.PathEscape(message))

	return doSendMessageRequest(messageURL, maasToken)
}

// doSendMessageRequest sends a request to the MaaS360 API to send a message to a device.
// It constructs the request, sends it, and processes the response.
func doSendMessageRequest(url string, maasToken string) error {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("MaaS token=\"%s\"", maasToken))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}
	var response actionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}
	if response.ActionResponse.Maas360DeviceID == "" {
		return fmt.Errorf("no device ID returned in response")
	}
	if response.ActionResponse.ActionStatus != 0 {
		return fmt.Errorf("message sending failed: %s", response.ActionResponse.Description)
	}

	fmt.Printf("Message sent successfully to device %s: %s\n", response.ActionResponse.Maas360DeviceID, response.ActionResponse.Description)
	return nil
}
