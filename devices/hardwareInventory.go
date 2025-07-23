package devices_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	auth_api "maas360api/auth"
	"maas360api/internal/constants"
)

type DeviceAttribute struct {
	AttributeKey   string `json:"key"`
	AttributeType  string `json:"type"`
	AttributeValue any    `json:"value"`
}

// Wrapper for the "deviceAttribute" slice
type DeviceAttributesWrapper struct {
	DeviceAttribute []DeviceAttribute `json:"deviceAttribute"`
}

type DeviceHardware struct {
	DeviceAttributes DeviceAttributesWrapper `json:"deviceAttributes"`
	ID               string                  `json:"maas360DeviceId"`
}

type HardwareInventoryResponse struct {
	DeviceHardware DeviceHardware `json:"deviceHardware"`
}

// GetHardwareInventory retrieves the hardware inventory for a specific device in MaaS360.
// It requires a billing ID, device ID, and an authentication token.
func GetHardwareInventory(billingID string, deviceID string, maasToken string) (*HardwareInventoryResponse, error) {
	if len(billingID) == 0 || len(deviceID) == 0 {
		return nil, fmt.Errorf("billing ID and device ID cannot be empty")
	}
	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}

	// Construct the hardware inventory URL
	hardwareInventoryURL := fmt.Sprintf("%s/device-apis/devices/1.0/hardwareInventory/%s?deviceId=%s", serviceURL, billingID, deviceID)

	// Perform the hardware inventory request
	return doHardwareInventoryRequest(hardwareInventoryURL, maasToken)
}

// doHardwareInventoryRequest sends a request to the MaaS360 API to retrieve hardware inventory for a device.
// It constructs the request, sends it, and processes the response.
func doHardwareInventoryRequest(url string, maasToken string) (*HardwareInventoryResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.AcceptHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.AuthorizationHeader, fmt.Sprintf(constants.MaaSTokenPrefix, maasToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var hardwareInventoryResp HardwareInventoryResponse
	if err := json.Unmarshal(body, &hardwareInventoryResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response JSON: %v; body: %s", err, string(body))
	}

	return &hardwareInventoryResp, nil
}

// PrintHardwareInventory prints the hardware inventory for a specific device in a human-readable format.
// It retrieves the hardware inventory using GetHardwareInventory and formats the output.
func PrintHardwareInventory(billingID string, deviceID string, maasToken string) {
	hardwareInventory, err := GetHardwareInventory(billingID, deviceID, maasToken)
	if err != nil {
		log.Fatalf("Error getting hardware inventory: %v", err)
	}
	if hardwareInventory == nil {
		log.Fatal("No hardware inventory found for the device")
	}
	fmt.Printf("Hardware Inventory for Device ID %s:\n", deviceID)
	for _, attr := range hardwareInventory.DeviceHardware.DeviceAttributes.DeviceAttribute {
		// Handle different types of attribute values
		if attr.AttributeValue == nil {
			fmt.Printf(" %s: <nil>\n", attr.AttributeKey)
			continue
		}
		switch v := attr.AttributeValue.(type) {
		case string:
			// Attempt to parse the string as time
			if t, err := time.Parse("2006-01-02T15:04:05", v); err == nil {
				fmt.Printf(" %s: %s\n", attr.AttributeKey, t.UTC().Format(time.RFC1123))
			} else {
				fmt.Printf(" %s: %s\n", attr.AttributeKey, v)
			}
		case float64:
			fmt.Printf(" %s: %.2f\n", attr.AttributeKey, v)
		case int:
			fmt.Printf(" %s: %d\n", attr.AttributeKey, v)
		case bool:
			fmt.Printf(" %s: %t\n", attr.AttributeKey, v)
		default:
			fmt.Printf(" %s: %v\n", attr.AttributeKey, v)
		}
	}
}
