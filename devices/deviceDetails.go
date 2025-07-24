package devices

import (
	"encoding/json"
	"fmt"
	"io"
	auth_api "maas360api/auth"
	httputil "maas360api/internal/http"
	"maas360api/internal/types"
)

type DeviceIdentifiers struct {
	Maas360DeviceID              string            `json:"maas360DeviceID"`
	DeviceName                   string            `json:"deviceName"`
	CustomAssetNumber            string            `json:"customAssetNumber"`
	Ownership                    string            `json:"ownership"`
	DeviceOwner                  string            `json:"deviceOwner"`
	Username                     string            `json:"username"`
	EmailAddress                 string            `json:"emailAddress"`
	PlatformName                 string            `json:"platformName"`
	SourceID                     int               `json:"sourceID"`
	DeviceType                   string            `json:"deviceType"`
	Manufacturer                 string            `json:"manufacturer"`
	Model                        string            `json:"model"`
	OSName                       string            `json:"osName"`
	OSServicePack                string            `json:"osServicePack"`
	IMEIESN                      any               `json:"imeiEsn"` // String or int64 | Empty string if not set
	InstalledDate                string            `json:"installedDate"`
	LastReported                 string            `json:"lastReported"`
	InstalledDateInEpochms       types.FlexibleInt `json:"installedDateInEpochms"`
	LastReportedInEpochms        types.FlexibleInt `json:"lastReportedInEpochms"`
	DeviceStatus                 string            `json:"deviceStatus"`
	Maas360ManagedStatus         string            `json:"maas360ManagedStatus"`
	UDID                         string            `json:"udid"`
	WifiMacAddress               string            `json:"wifiMacAddress"`
	MailboxDeviceId              string            `json:"mailboxDeviceId"`
	MailboxLastReported          string            `json:"mailboxLastReported"`
	MailboxLastReportedInEpochms types.FlexibleInt `json:"mailboxLastReportedInEpochms"`
	MailboxManaged               string            `json:"mailboxManaged"`
	IsSupervisedDevice           bool              `json:"isSupervisedDevice"`
	TestDevice                   bool              `json:"testDevice"`
	UnifiedTravelerDeviceId      string            `json:"unifiedTravelerDeviceId"`
}

type DeviceResponse struct {
	Device DeviceIdentifiers `json:"device"`
}

func GetDevice(billingID string, deviceID string, maasToken string) (*DeviceIdentifiers, error) {
	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}

	searchURL := fmt.Sprintf("%s/device-apis/devices/1.0/core/%s?deviceId=%s", serviceURL, billingID, deviceID)

	resp, err := httputil.DoMaaSRequest(httputil.RequestOptions{
		Method:    "GET",
		URL:       searchURL,
		MaaSToken: maasToken,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
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
