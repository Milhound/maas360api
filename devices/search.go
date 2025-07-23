package devices_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	auth_api "maas360api/auth"
)

// Device represents a MaaS360 device.
type Device struct {
	AppComplianceStatus          string `json:"appComplianceState"`
	AssetTag                     string `json:"customAssetNumber"`
	DeviceOwner                  string `json:"deviceOwner"`
	Email                        string `json:"emailAddress"`
	EncryptionStatus             string `json:"encryptionStatus"`
	FirstRegisteredInEpochms     any    `json:"firstRegisteredInEpochms"` // String or int64 | Empty string if not set
	ID                           any    `json:"maas360DeviceID"`          // String or int32 | Empty string if not set
	IMEI                         any    `json:"imeiEsn"`                  // String or int64 | Empty string if not set
	InstalledDate                string `json:"installedDate"`
	InstalledDateInEpochms       any    `json:"installedDateInEpochms"` // String or int64 | Empty string if not set
	JailbreakStatus              string `json:"jailbreakStatus"`
	LastMDMRegisteredInEpochms   any    `json:"lastMdmRegisteredInEpochms"` // String or int64 | Empty string if not set
	LastReported                 string `json:"lastReported"`
	LastReportedInEpochms        any    `json:"lastReportedInEpochms"`   // String or int64 | Empty string if not set
	LastRegisteredInEpochms      any    `json:"lastRegisteredInEpochms"` // String or int64 | Empty string if not set
	MacAddress                   string `json:"wifiMacAddress"`
	MailboxID                    any    `json:"mailboxDeviceId"` // String or int64 | Empty string if not set
	MailboxLastRepoted           string `json:"mailboxLastReported"`
	MailboxLastReportedInEpochms any    `json:"mailboxLastReportedInEpochms"` // String or int64 | Empty string if not set
	MailboxStatus                string `json:"mailboxManaged"`
	MDMPolicy                    string `json:"mdmPolicy"`
	MDMMailboxDeviceID           string `json:"mdmMailboxDeviceId"`
	ManagedStatus                string `json:"maas360ManagedStatus"`
	Manufacturer                 string `json:"manufacturer"`
	Model                        string `json:"modelId"`
	Name                         string `json:"deviceName"`
	OS                           string `json:"osName"`
	OSVersion                    any    `json:"osVersion"` // String or int32
	Ownership                    string `json:"ownership"`
	PasscodeComplianceStatus     string `json:"passcodeCompliance"`
	PhoneNumber                  any    `json:"phoneNumber"` // String or int64 | Empty string if not set
	Platform                     string `json:"platformName"`
	PolicyComplianceStatus       string `json:"policyComplianceState"`
	RuleComplianceStatus         string `json:"ruleComplianceState"`
	SelectiveWipeStatus          string `json:"selectiveWipeStatus"`
	SerialNumber                 string `json:"platformSerialNumber"`
	ServicePack                  string `json:"osServicePack"`
	Status                       string `json:"deviceStatus"`
	Supervised                   bool   `json:"isSupervisedDevice"`
	TestDevice                   bool   `json:"testDevice"`
	TravelerDeviceID             any    `json:"unifiedTravelerDeviceId"` // String or int64 | Empty string if not set
	Type                         string `json:"deviceType"`
	UDID                         string `json:"udid"`
	UserDomain                   string `json:"userDomain"`
	Username                     string `json:"username"`
	EnrollmentMode               string `json:"enrollmentMode"`
	SourceID                     int32  `json:"sourceID"`
}

type DeviceOrDevices []Device

type devices struct {
	Device DeviceOrDevices `json:"device"`
}

type searchResponse struct {
	Devices devices `json:"devices"`
}

func (d *DeviceOrDevices) UnmarshalJSON(data []byte) error {
	// Try as array
	var arr []Device
	if err := json.Unmarshal(data, &arr); err == nil {
		*d = arr
		return nil
	}
	// Try as single object
	var single Device
	if err := json.Unmarshal(data, &single); err == nil {
		*d = []Device{single}
		return nil
	}
	return fmt.Errorf("DeviceOrDevices: cannot unmarshal %s", string(data))
}

// Search performs a search for devices in the MaaS360 API based on the provided filters.
// It returns a list of devices that match the search criteria.
func SearchDevices(billingID string, token string, filters map[string]string) ([]Device, error) {
	// Possible search filters:
	// "deviceStatus": "InActive", // ["Active", "InActive"] Default is "Active"
	// "partialDeviceName":   "",
	// "partialUsername": "",
	// "partialPhoneNumber": "",
	// "udid": "",
	// "imeiMeid": "",
	// "wifiMacAddress": "",
	// "mailboxDeviceId": "",
	// "platformName": "iOS", // ["iOS", "Android", "Windows", "Mac", "Others"]
	// "excludeCloudExtenders": "No", // Default is "Yes"
	// "maas360DeviceId ": "",
	// "userDomain": "",
	// "email": "",
	// "maas360ManagedStatus ": "Activated", // ["Inactive", "Activated", "Control Removed", "Pending Control Removed", "User Removed Control", "Not Enrolled", "Enrolled"]
	// "mailBoxManaged": "", // ["ActiveSync", "Domino", "BES", "GmailSync"]
	// "mdmMailboxDeviceId": "",
	// "plcCompliance": "OOC", // ["OOC", "ALL"] Default is "ALL"
	// "ruleCompliance": "OOC", // ["OOC", "ALL"] Default is "ALL"
	// "appCompliance": "OOC", // ["OOC", "ALL"] Default is "ALL"
	// "pswdCompliance": "OOC", // ["OOC", "ALL"] Default is "ALL"

	// Validate required fields
	if len(billingID) == 0 || len(token) == 0 {
		return nil, fmt.Errorf("billingID and token are required")
	}

	// Get the MaaS360 instance URL
	instance, err := auth_api.GetInstance(billingID)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(filters)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}
	searchFilters := url.Values{}
	if len(jsonData) > 0 {
		for key, value := range filters {
			searchFilters.Add(key, value)
		}
	}

	searchURL := fmt.Sprintf("%s/device-apis/devices/2.0/search/customer/%s?", instance, billingID) + searchFilters.Encode()

	return doSearchDevicesRequest(searchURL, token)
}

// doSearchRequest sends a search request to the MaaS360 API and returns the list of devices.
// It constructs the request, sends it, and processes the response.
func doSearchDevicesRequest(url string, token string) ([]Device, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
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
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	var devicesResp searchResponse
	if err := json.Unmarshal(body, &devicesResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response JSON: %v; body: %s", err, string(body))
	}
	if len(devicesResp.Devices.Device) == 0 {
		return nil, fmt.Errorf("no devices found")
	}
	return devicesResp.Devices.Device, nil
}

// PrintDevices retrieves and prints the list of devices based on the provided filters.
// It calls SearchDevices to get the devices and then logs their details.
func PrintDevices(billingID string, token string, filters map[string]string) {
	devices, err := SearchDevices(billingID, token, filters)
	if err != nil {
		log.Fatalf("Error searching devices: %v", err)
	}
	if len(devices) == 0 {
		log.Fatal("No devices found")
	}
	log.Printf("Found %d devices", len(devices))
	for _, device := range devices {
		log.Printf("Device Name: %s, CSN: %s, Status: %s", device.Name, device.ID, device.Status)
	}
}
