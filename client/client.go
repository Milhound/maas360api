package client

import (
	"errors"
	"time"

	"maas360api/application"
	"maas360api/auth"
	"maas360api/devices"
	"maas360api/internal/constants"
)

// MaaS360Client represents a MaaS360 API client with authentication credentials
// and methods for interacting with the MaaS360 API.
type MaaS360Client struct {
	BillingID    string // MaaS360 billing ID
	AppID        string // Application ID for API access
	AccessKey    string // Application access key
	Username     string // Username for authentication
	Password     string // Password for authentication
	RefreshToken string // Refresh token for token-based authentication
	MaasToken    string // Current authentication token
}

// GetBasicauth generates a Basic Authentication header value for the client's credentials.
// Returns an empty string if username or password is empty.
func (c *MaaS360Client) GetBasicauth() string {
	return auth.GetBasicAuth(c.Username, c.Password)
}

// Authenticate obtains an authentication token from the MaaS360 API.
// This must be called before using other API methods.
// Returns an error if authentication fails or if no valid token is received.
func Authenticate(credentials auth.MaaS360AdminAuth) (*MaaS360Client, error) {
	credentials.PlatformID = constants.Platform
	credentials.AppVersion = constants.Version

	authResponse, err := auth.Auth(credentials)
	if err != nil {
		return nil, err
	}
	if authResponse.AuthToken == "" {
		return nil, errors.New("failed to retrieve MaaS360 auth token")
	}
	if authResponse.RefreshToken == "" {
		return nil, errors.New("failed to retrieve MaaS360 refresh token")
	}

	c := &MaaS360Client{
		BillingID:    credentials.BillingID,
		AppID:        credentials.AppID,
		AccessKey:    credentials.AccessKey,
		Username:     credentials.Username,
		Password:     credentials.Password,
		RefreshToken: authResponse.RefreshToken,
		MaasToken:    authResponse.AuthToken,
	}
	return c, nil
}
func (c *MaaS360Client) GetDeviceActions(deviceID string) (*devices.DeviceActionsResponse, error) {
	return devices.GetDeviceActions(c.BillingID, deviceID, c.MaasToken)
}

func (c *MaaS360Client) PerformDeviceAction(deviceID string, actionID string, additionalParams map[string]string) error {
	return devices.PerformDeviceAction(c.BillingID, deviceID, actionID, additionalParams, c.MaasToken)
}

func (c *MaaS360Client) SendMessage(deviceID string, subject string, message string) error {
	return devices.SendMessage(c.BillingID, deviceID, subject, message, c.MaasToken)
}

func (c *MaaS360Client) LockDevice(deviceID string) error {
	return devices.LockDevice(c.BillingID, deviceID, c.MaasToken)
}

func (c *MaaS360Client) GetHardwareInventory(deviceID string) (*devices.HardwareInventoryResponse, error) {
	return devices.GetHardwareInventory(c.BillingID, deviceID, c.MaasToken)
}

func (c *MaaS360Client) PrintHardwareInventory(deviceID string) {
	devices.PrintHardwareInventory(c.BillingID, deviceID, c.MaasToken)
}

func (c *MaaS360Client) GetSoftwareInstalled(deviceID string) (*devices.SoftwareInstalledResponse, error) {
	return devices.GetSoftwareInstalled(c.BillingID, deviceID, c.MaasToken)
}

func (c *MaaS360Client) PrintSoftwareInstalled(deviceID string) {
	devices.PrintSoftwareInstalled(c.BillingID, deviceID, c.MaasToken)
}

func (c *MaaS360Client) GetDevice(deviceID string) (*devices.DeviceIdentifiers, error) {
	return devices.GetDevice(c.BillingID, deviceID, c.MaasToken)
}

func (c *MaaS360Client) SearchDevices(filters map[string]string) ([]devices.Device, error) {
	return devices.SearchDevices(c.BillingID, filters, c.MaasToken)
}

func (c *MaaS360Client) PrintDevices(filters map[string]string) {
	devices.PrintDevices(c.BillingID, filters, c.MaasToken)
}

func (c *MaaS360Client) SearchCatalog(filters map[string]string) ([]application.CatalogApp, error) {
	return application.SearchCatalog(c.BillingID, filters, c.MaasToken)
}

func (c *MaaS360Client) PrintCatalogApps(filters map[string]string) {
	application.PrintCatalogApps(c.BillingID, filters, c.MaasToken)
}

func (c *MaaS360Client) SearchInstalledApps(filters map[string]string) ([]application.InstalledApp, error) {
	return application.SearchInstalledApps(c.BillingID, filters, c.MaasToken)
}

func (c *MaaS360Client) PrintAllSoftwareInstalled(filters map[string]string) {
	application.PrintAllSoftwareInstalled(c.BillingID, filters, c.MaasToken)
}

func (c *MaaS360Client) UpdateOS(deviceID string, osVersion string, targetLocalTime time.Time) error {
	return devices.UpdateOS(c.BillingID, deviceID, osVersion, targetLocalTime, c.MaasToken)
}

func (c *MaaS360Client) GetNetworkInfo(deviceID string) ([]devices.DeviceAttribute, error) {
	return devices.GetNetworkInfo(c.BillingID, deviceID, c.MaasToken)
}
func (c *MaaS360Client) PrintNetworkInfo(deviceID string) {
	devices.PrintNetworkInfo(c.BillingID, deviceID, c.MaasToken)
}
