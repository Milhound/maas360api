package devices_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	auth_api "maas360api/auth"
)

type Attribute struct {
	AttributeKey   string `json:"key"`
	AttributeType  string `json:"type"`
	AttributeValue any    `json:"value"` // changed to `any` to handle string, float, int
}

type Software struct {
	Name       string      `json:"swName"`
	Attributes []Attribute `json:"swAttrs"`
}

type DeviceSoftwares struct {
	ID                  string     `json:"maas360DeviceID"`
	LastDataRefreshTime string     `json:"lastSoftwareDataRefreshDate"`
	Softwares           []Software `json:"deviceSw"` // Correct key and case
}

type SoftwareInstalledResponse struct {
	DeviceSoftwares DeviceSoftwares `json:"deviceSoftwares"`
}

// GetSoftwareInstalled retrieves the software installed on a specific device in MaaS360.
func GetSoftwareInstalled(billingID string, deviceID string, maasToken string) (*SoftwareInstalledResponse, error) {
	if len(billingID) == 0 || len(deviceID) == 0 {
		return nil, fmt.Errorf("billing ID and device ID cannot be empty")
	}
	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}

	// Construct the hardware inventory URL
	softwareInstalledURL := fmt.Sprintf("%s/device-apis/devices/1.0/softwareInstalled/%s?deviceId=%s", serviceURL, billingID, deviceID)

	// Perform the hardware inventory request
	return doGetSoftwareInstalled(softwareInstalledURL, maasToken)
}

// doGetSoftwareInstalled sends a request to the MaaS360 API to retrieve software installed on a device.
// It constructs the request, sends it, and processes the response.
func doGetSoftwareInstalled(url string, maasToken string) (*SoftwareInstalledResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("MaaS token=\"%s\"", maasToken))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

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

	var softwareInstalledResponse SoftwareInstalledResponse
	if err := json.Unmarshal(body, &softwareInstalledResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response JSON: %v; body: %s", err, string(body))
	}
	if softwareInstalledResponse.DeviceSoftwares.ID == "" {
		return nil, fmt.Errorf("no software installed data found for device")
	}
	return &softwareInstalledResponse, nil
}

// PrintSoftwareInstalled prints the software installed on a specific device in a human-readable format.
// It retrieves the software installed using GetSoftwareInstalled and formats the output.
func PrintSoftwareInstalled(billingID string, deviceID string, maasToken string) {
	softwareInstalled, err := GetSoftwareInstalled(billingID, deviceID, maasToken)
	if err != nil {
		log.Fatalf("Error getting hardware inventory: %v", err)
	}
	if softwareInstalled == nil {
		log.Fatal("No hardware inventory found for the device")
	}
	fmt.Printf("Software Installed for Device ID %s:\n", deviceID)
	fmt.Printf("Last Data Refresh Time: %s\n", softwareInstalled.DeviceSoftwares.LastDataRefreshTime)
	for _, attr := range softwareInstalled.DeviceSoftwares.Softwares {
		fmt.Printf("Software Name: %s\n", attr.Name)
		for _, attribute := range attr.Attributes {
			// Handle different types of attribute values
			if attribute.AttributeValue == nil {
				fmt.Printf(" %s: <nil>\n", attribute.AttributeKey)
				continue
			}
			switch v := attribute.AttributeValue.(type) {
			case string:
				fmt.Printf(" %s: %s\n", attribute.AttributeKey, v)
			case float64:
				fmt.Printf(" %s: %.2f\n", attribute.AttributeKey, v)
			case int:
				fmt.Printf(" %s: %d\n", attribute.AttributeKey, v)
			default:
				fmt.Printf(" %s: %v\n", attribute.AttributeKey, v)
			}
		}
	}
}
