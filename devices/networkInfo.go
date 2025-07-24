package devices

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maas360api/auth"
	httputil "maas360api/internal/http"
	"net/http"
)

type NetworkInformationWrapper struct {
	NetworkInformation DeviceAttributesResponse `json:"networkInformation"`
}

func GetNetworkInfo(billingID string, deviceID string, maasToken string) ([]DeviceAttribute, error) {
	if billingID == "" || deviceID == "" || maasToken == "" {
		return nil, fmt.Errorf("billingID, deviceID, and maasToken must not be empty")
	}

	serviceURL, err := auth.GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}

	searchURL := fmt.Sprintf("%s/device-apis/devices/1.0/mdNetworkInformation/%s?deviceId=%s", serviceURL, billingID, deviceID)

	resp, err := httputil.DoMaaSRequest(httputil.RequestOptions{
		Method:    "GET",
		URL:       searchURL,
		MaaSToken: maasToken,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get network info: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wrapper NetworkInformationWrapper
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	fmt.Println(string(body))
	return wrapper.NetworkInformation.AttributeWrapper.DeviceAttributes, nil
}

func PrintNetworkInfo(billingID string, deviceID string, maasToken string) {
	networkInfo, err := GetNetworkInfo(billingID, deviceID, maasToken)
	if err != nil {
		log.Fatalf("Error getting network info: %v", err)
	}
	if networkInfo == nil {
		log.Fatal("No network info found for the device")
	}
	fmt.Printf("Network Info for Device ID %s:\n", deviceID)
	for _, attr := range networkInfo {
		fmt.Printf(" %s: %v\n", attr.AttributeKey, attr.AttributeValue)
	}
}
