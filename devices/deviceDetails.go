package devices_api

import (
	"encoding/json"
	"fmt"
	"io"
	auth_api "maas360api/auth"
	"net/http"
)

type CoreAttributes struct {
	Maas360DeviceID              string `json:"maas360DeviceID"`
	DeviceName                   string `json:"deviceName"`
	CustomAssetNumber            string `json:"customAssetNumber"`
	Ownership                    string `json:"ownership"`
	DeviceOwner                  string `json:"deviceOwner"`
	Username                     string `json:"username"`
	EmailAddress                 string `json:"emailAddress"`
	PlatformName                 string `json:"platformName"`
	SourceID                     int    `json:"sourceID"`
	DeviceType                   string `json:"deviceType"`
	Manufacturer                 string `json:"manufacturer"`
	Model                        string `json:"model"`
	OSName                       string `json:"osName"`
	OSServicePack                string `json:"osServicePack"`
	IMEIESN                      any    `json:"imeiEsn"` // String or int64 | Empty string if not set
	InstalledDate                string `json:"installedDate"`
	LastReported                 string `json:"lastReported"`
	InstalledDateInEpochms       any    `json:"installedDateInEpochms"` // String or int64 | Empty string if not set
	LastReportedInEpochms        any    `json:"lastReportedInEpochms"`  // String or int64 | Empty string if not set
	DeviceStatus                 string `json:"deviceStatus"`
	Maas360ManagedStatus         string `json:"maas360ManagedStatus"`
	UDID                         string `json:"udid"`
	WifiMacAddress               string `json:"wifiMacAddress"`
	MailboxDeviceId              string `json:"mailboxDeviceId"`
	MailboxLastReported          string `json:"mailboxLastReported"`
	MailboxLastReportedInEpochms any    `json:"mailboxLastReportedInEpochms"` // String or int64 | Empty string if not set
	MailboxManaged               string `json:"mailboxManaged"`
	IsSupervisedDevice           bool   `json:"isSupervisedDevice"`
	TestDevice                   bool   `json:"testDevice"`
	UnifiedTravelerDeviceId      string `json:"unifiedTravelerDeviceId"`
}

type DeviceResponse struct {
	Device CoreAttributes `json:"device"`
}

func GetDevice(billingID string, deviceID string, maasToken string) (*CoreAttributes, error) {
	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}

	searchURL := fmt.Sprintf("%s/device-apis/devices/1.0/core/%s?deviceId=%s", serviceURL, billingID, deviceID)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("MaaS token=\"%s\"", maasToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var deviceResp DeviceResponse
	if err := json.Unmarshal(body, &deviceResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response JSON: %v; body: %s", err, string(body))
	}

	return &deviceResp.Device, nil
}
