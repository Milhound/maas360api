package client

import (
	"log"

	application_api "maas360api/application"
	auth_api "maas360api/auth"
	device_api "maas360api/devices"
)

type Client struct {
	BillingID string
	AppID     string
	AccessKey string
	Username  string
	Password  string
	Refresh   string
	MaasToken string
}

func NewClient(billingID, appID, accessKey, username, password, refresh string) *Client {
	return &Client{
		BillingID: billingID,
		AppID:     appID,
		AccessKey: accessKey,
		Username:  username,
		Password:  password,
		Refresh:   refresh,
	}
}

func (c *Client) GetBasicauth() string {
	return auth_api.GetBasicAuth(c.Username, c.Password)
}

func (c *Client) Authenticate() error {
	var err error
	c.MaasToken, _, err = auth_api.GetToken(c.BillingID, c.AppID, c.AccessKey, c.Username, c.Password, c.Refresh)
	if err != nil {
		return err
	}
	if c.MaasToken == "" {
		return log.Output(2, "Failed to retrieve MaaS360 auth token")
	}
	return nil
}
func (c *Client) GetDeviceActions(deviceID string) (*device_api.DeviceActionsResponse, error) {
	return device_api.GetDeviceActions(c.BillingID, deviceID, c.MaasToken)
}

func (c *Client) PerformDeviceAction(deviceID string, actionID string, additionalParams map[string]string) error {
	return device_api.PerformDeviceAction(c.BillingID, deviceID, actionID, additionalParams, c.MaasToken)
}

func (c *Client) SendMessage(deviceID string, subject string, message string) error {
	return device_api.SendMessage(c.BillingID, deviceID, subject, message, c.MaasToken)
}

func (c *Client) LockDevice(deviceID string) error {
	return device_api.LockDevice(c.BillingID, deviceID, c.MaasToken)
}

func (c *Client) GetHardwareInventory(deviceID string) (*device_api.HardwareInventoryResponse, error) {
	return device_api.GetHardwareInventory(c.BillingID, deviceID, c.MaasToken)
}

func (c *Client) PrintHardwareInventory(deviceID string) {
	device_api.PrintHardwareInventory(c.BillingID, deviceID, c.MaasToken)
}

func (c *Client) GetSoftwareInstalled(deviceID string) (*device_api.SoftwareInstalledResponse, error) {
	return device_api.GetSoftwareInstalled(c.BillingID, deviceID, c.MaasToken)
}

func (c *Client) PrintSoftwareInstalled(deviceID string) {
	device_api.PrintSoftwareInstalled(c.BillingID, deviceID, c.MaasToken)
}

func (c *Client) GetDevice(deviceID string) (*device_api.CoreAttributes, error) {
	return device_api.GetDevice(c.BillingID, deviceID, c.MaasToken)
}

func (c *Client) SearchDevices(filters map[string]string) ([]device_api.Device, error) {
	return device_api.SearchDevices(c.BillingID, filters, c.MaasToken)
}

func (c *Client) PrintDevices(filters map[string]string) {
	device_api.PrintDevices(c.BillingID, filters, c.MaasToken)
}

func (c *Client) SearchCatalog(filters map[string]string) ([]application_api.CatalogApp, error) {
	return application_api.SearchCatalog(c.BillingID, filters, c.MaasToken)
}

func (c *Client) PrintCatalogApps(filters map[string]string) {
	application_api.PrintCatalogApps(c.BillingID, filters, c.MaasToken)
}

func (c *Client) SearchInstalledApps(filters map[string]string) ([]application_api.InstalledApp, error) {
	return application_api.SearchInstalledApps(c.BillingID, filters, c.MaasToken)
}

func (c *Client) PrintAllSoftwareInstalled(filters map[string]string) {
	application_api.PrintAllSoftwareInstalled(c.BillingID, filters, c.MaasToken)
}
