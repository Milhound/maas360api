package devices_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	auth_api "maas360api/auth"
)

type DeviceAction struct {
	ActionID    string `json:"actionId"`
	ActionName  string `json:"actionName"`
	ActionOrder int    `json:"actionOrder"`
	ActionType  string `json:"actionType"`
}

type DeviceActions struct {
	Actions []DeviceAction `json:"deviceAction"`
}

type DeviceActionsResponse struct {
	DeviceActions DeviceActions `json:"deviceActions"`
}

type ActionRequest struct {
	Name             string            `json:"name"`
	ExpiryDate       int64             `json:"expiryDate"`
	RequesterWorflow string            `json:"requesterWorkflow"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"`
}

func (d *DeviceActionsResponse) GetActionByName(actionName string) (*DeviceAction, error) {
	for _, action := range d.DeviceActions.Actions {
		if action.ActionName == actionName {
			return &action, nil
		}
	}
	return nil, fmt.Errorf("action with name %s not found", actionName)
}

func (d *DeviceActionsResponse) GetActionByID(actionID string) (*DeviceAction, error) {
	for _, action := range d.DeviceActions.Actions {
		if action.ActionID == actionID {
			return &action, nil
		}
	}
	return nil, fmt.Errorf("action with ID %s not found", actionID)
}

// GetDeviceActions retrieves the list of available device actions for a specific device.
func GetDeviceActions(billingID string, deviceID string, token string) (*DeviceActionsResponse, error) {
	if billingID == "" || deviceID == "" || token == "" {
		return nil, fmt.Errorf("billingID, deviceID, and token must not be empty")
	}
	instance, err := auth_api.GetInstance(billingID)
	if err != nil {
		return nil, fmt.Errorf("error getting instance: %v", err)
	}

	url := fmt.Sprintf("%s/device-apis/devices/1.0/deviceActions/%s?deviceId=%s", instance, billingID, deviceID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("MaaS token=\"%s\"", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var response DeviceActionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response body: %s\n", string(bodyBytes))
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &response, nil
}

// PerformDeviceAction performs a specific action on a device.
func PerformDeviceAction(billingID string, deviceID string, actionID string, additionalParams map[string]string, token string, basicAuth string) error {
	// Notes:
	// MDM_LOCATE: Not applicable action for iOS devices
	//

	if billingID == "" || deviceID == "" || actionID == "" || token == "" || basicAuth == "" {
		return fmt.Errorf("billingID, deviceID, actionID, and basicAuth must not be empty")
	}

	actionsResponse, err := GetDeviceActions(billingID, deviceID, token)
	if err != nil {
		return fmt.Errorf("error getting device actions: %v", err)
	}

	action, err := actionsResponse.GetActionByID(actionID)
	if err != nil {
		return fmt.Errorf("error getting action by name: %v", err)
	}

	if actionID == "ANDROID_CUSTOM_CMDS" && additionalParams == nil {
		return fmt.Errorf("additionalParams must not be nil for action %s", actionID)
	}

	fmt.Printf("Performing action: %s\n", action.ActionName)
	if actionID == "ANDROID_CUSTOM_CMDS" && additionalParams != nil {
		err = doAction(billingID, deviceID, action.ActionID, action.ActionName, additionalParams, basicAuth)
	} else {
		err = doAction(billingID, deviceID, action.ActionID, action.ActionName, nil, basicAuth)
	}

	if err != nil {
		return fmt.Errorf("error performing action: %v", err)
	}
	return nil
}

// doAction sends a request to perform a specific action on a device.
// It constructs the request, sends it, and processes the response.
func doAction(billingID string, deviceID string, actionID string, actionName string, additionalParams map[string]string, basicAuth string) error {
	if billingID == "" || deviceID == "" || actionID == "" || actionName == "" || basicAuth == "" {
		return fmt.Errorf("billingID, deviceID, actionName, and basicAuth must not be empty")
	}
	instance, err := auth_api.GetInstance(billingID)
	if err != nil {
		return fmt.Errorf("error getting instance: %v", err)
	}
	var reqBodyRaw ActionRequest
	reqBodyRaw.Name = actionName
	reqBodyRaw.ExpiryDate = time.Now().Local().Unix() + 300 // 5 minutes from now
	reqBodyRaw.RequesterWorflow = "TEST"

	if actionID == "ANDROID_CUSTOM_CMDS" && additionalParams == nil {
		return fmt.Errorf("additionalParams must not be nil for action %s", actionID)
	}

	if actionID == "ANDROID_CUSTOM_CMDS" && additionalParams != nil {
		reqBodyRaw.AdditionalParams = additionalParams
	}

	reqBody, err := json.Marshal(reqBodyRaw)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	url := fmt.Sprintf("%s/action-apis/actions/1.0/customer/%s/action/%s/device/%s", instance, billingID, actionID, deviceID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", basicAuth)
	//req.Header.Set("Authorization", fmt.Sprintf("MaaS token=\"%s\"", token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request URL:", url)
		fmt.Println("Request Headers:")
		for name, values := range req.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}
		fmt.Println("Request Body:", string(reqBody))
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	fmt.Println("Action performed successfully")

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response body: %s\n", string(bodyBytes))

	return nil
}
