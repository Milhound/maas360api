package auth_api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"maas360api/internal/constants"
	httputil "maas360api/internal/http"
)

const (
	M1 = "https://services.fiberlink.com"
	M2 = "https://services.m2.maas360.com"
	M3 = "https://services.m3.maas360.com"
	M4 = "https://services.m4.maas360.com"
	M6 = "https://services.m6.maas360.com"
)

type MaaS360AdminAuth struct {
	BillingID    string `json:"billingID"`
	PlatformID   string `json:"platformID"`
	AppVersion   string `json:"appVersion"`
	AppID        string `json:"appID"`
	AccessKey    string `json:"appAccessKey"`
	Username     string `json:"userName"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type authRequest struct {
	Auth MaaS360AdminAuth `json:"maaS360AdminAuth"`
}

type authParams struct {
	Request authRequest `json:"authRequest"`
}

type AuthResponseBody struct {
	ErrorCode    uint16 `json:"errorCode"`
	ErrorDesc    string `json:"errorDesc"`
	AuthToken    string `json:"authToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type AuthResponse struct {
	Wrapper AuthResponseBody `json:"authResponse"`
}

// client is the HTTP client used for making requests to the MaaS360 API.
var client = httputil.GetSharedClient()

// Authenticate with username/password api call
// This function sends a request to the MaaS360 authentication API to get an auth token.
func Auth(authCredentials MaaS360AdminAuth) (*AuthResponseBody, error) {
	if authCredentials.BillingID == "" || authCredentials.AppID == "" || authCredentials.AccessKey == "" || authCredentials.Username == "" {
		return nil, fmt.Errorf("billingID, appID, accessKey, and username must not be empty")
	}
	if authCredentials.Password == "" && authCredentials.RefreshToken == "" {
		return nil, fmt.Errorf("either password or refresh token must be provided for authentication")
	}
	if authCredentials.PlatformID == "" {
		authCredentials.PlatformID = constants.Platform
	}
	if authCredentials.AppVersion == "" {
		authCredentials.AppVersion = constants.Version
	}

	params := authParams{
		Request: authRequest{
			Auth: authCredentials,
		},
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}
	serviceURL, err := GetServiceURL(authCredentials.BillingID)
	if err != nil {
		return nil, fmt.Errorf("error getting service URL: %v", err)
	}
	url := fmt.Sprintf("%s%s/customer/%s", serviceURL, "/auth-apis/auth/2.0/authenticate", authCredentials.BillingID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.AcceptHeader, constants.ContentTypeJSON)

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

	var parsed AuthResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("error unmarshaling response JSON: %v; body: %s", err, string(body))
	}
	if parsed.Wrapper.ErrorCode != 0 {
		return nil, fmt.Errorf("error from MaaS360: %s (code: %d)", parsed.Wrapper.ErrorDesc, parsed.Wrapper.ErrorCode)
	}

	return &parsed.Wrapper, nil
}

// GetServiceURL returns the MaaS360 service URL based on the billing ID.
// It checks the first character of the billing ID to determine the instance.
func GetServiceURL(billingID string) (string, error) {
	if len(billingID) == 0 {
		return "", fmt.Errorf("billing ID cannot be empty")
	}
	switch billingID[0] {
	case '1':
		return M1, nil
	case '2':
		return M2, nil
	case '3':
		return M3, nil
	case '4':
		return M4, nil
	case '6':
		return M6, nil
	default:
		return "", fmt.Errorf("invalid billing ID: %s", billingID)
	}
}

// GetBasicAuth returns a Basic Auth header value for the given username and password. Including "Basic " prefix.
func GetBasicAuth(username, password string) string {
	if username == "" || password == "" {
		return ""
	}
	rawText := fmt.Sprintf("%s:%s", username, password)
	encodedText := base64.StdEncoding.EncodeToString([]byte(rawText))
	return fmt.Sprintf("Basic %s", encodedText)
}
