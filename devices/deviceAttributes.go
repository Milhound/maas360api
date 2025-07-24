package devices

import (
	"encoding/json"
	"fmt"
	"io"
	httputil "maas360api/internal/http"
	"net/http"
)

type CustomAttribute struct {
	Name  string `json:"customAttributeName"`
	Value any    `json:"customAttributeValue"`
}

type CustomAttributesWrapper struct {
	CustomAttribute []CustomAttribute `json:"customAttribute"`
}

type DeviceIdentity struct {
	DeviceID               string                  `json:"maas360DeviceID"`
	CustomAssetNumber      string                  `json:"customAssetNumber"`
	Owner                  string                  `json:"owner"`
	Ownership              string                  `json:"ownership"`
	Vendor                 string                  `json:"vendor"`
	PoNumber               string                  `json:"poNumber"`
	PurchaseType           string                  `json:"purchaseType"`
	PurchaseDate           string                  `json:"purchaseDate"`
	PurchasePrice          string                  `json:"purchasePrice"`
	WarrantyNumber         string                  `json:"warrantyNumber"`
	WarrantyExpirationDate string                  `json:"warrantyExpirationDate"`
	WarrantyType           string                  `json:"warrantyType"`
	Office                 string                  `json:"office,omitempty"`
	Department             string                  `json:"department"`
	CustomAttributes       CustomAttributesWrapper `json:"customAttributes"`
}

type DeviceIdentityResponse struct {
	DeviceIdentity DeviceIdentity `json:"deviceIdentity"`
}

func GetDeviceAttributes(serviceURL string, billingID string, deviceID string, maasToken string) (*DeviceIdentity, error) {
	if serviceURL == "" || billingID == "" || deviceID == "" || maasToken == "" {
		return nil, fmt.Errorf("serviceURL, billingID, deviceID, and maasToken must not be empty")
	}

	url := fmt.Sprintf("%s/device-apis/devices/1.0/identity/%s?deviceId=%s", serviceURL, billingID, deviceID)

	resp, err := httputil.DoMaaSRequest(httputil.RequestOptions{
		Method:    "GET",
		URL:       url,
		MaaSToken: maasToken,
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var attributesResponse DeviceIdentityResponse
	if err := json.Unmarshal(body, &attributesResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &attributesResponse.DeviceIdentity, nil
}

func PrintDeviceAttributes(serviceURL string, billingID string, deviceID string, maasToken string) {
	identity, err := GetDeviceAttributes(serviceURL, billingID, deviceID, maasToken)
	if err != nil {
		fmt.Printf("Error retrieving device attributes: %v\n", err)
		return
	}

	fmt.Printf("Device Attributes for Device ID %s:\n", deviceID)
	fmt.Printf(" Ownership: %s\n", identity.Ownership)
	fmt.Printf(" Office: %s\n", identity.Office)
	fmt.Printf(" Department: %s\n", identity.Department)
	fmt.Printf(" Vendor: %s\n", identity.Vendor)
	fmt.Printf(" PO Number: %s\n", identity.PoNumber)
	fmt.Printf(" Purchase Type: %s\n", identity.PurchaseType)
	fmt.Printf(" Purchase Date: %s\n", identity.PurchaseDate)
	fmt.Printf(" Purchase Price: %s\n", identity.PurchasePrice)
	fmt.Printf(" Warranty Number: %s\n", identity.WarrantyNumber)
	fmt.Printf(" Warranty Expiration Date: %s\n", identity.WarrantyExpirationDate)
	fmt.Printf(" Warranty Type: %s\n", identity.WarrantyType)
	fmt.Printf(" Custom Asset Number: %s\n", identity.CustomAssetNumber)
	fmt.Printf(" Owner: %s\n", identity.Owner)
	fmt.Println(" Custom Attributes:")
	// Print custom attributes
	for _, attr := range identity.CustomAttributes.CustomAttribute {
		fmt.Printf(" - %s: %v\n", attr.Name, attr.Value)
	}
}
