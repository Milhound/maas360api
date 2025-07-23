package auth_api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

type maaS360AdminAuth struct {
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
	Auth maaS360AdminAuth `json:"maaS360AdminAuth"`
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

// helper function to do the auth or refresh call
// It constructs the request, sends it, and processes the response.
func doAuthRequest(serviceURL string, billingID string, path string, auth maaS360AdminAuth, extraHeaders map[string]string) (*AuthResponse, error) {
	params := authParams{
		Request: authRequest{
			Auth: auth,
		},
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	url := fmt.Sprintf("%s%s/customer/%s", serviceURL, path, billingID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.AcceptHeader, constants.ContentTypeJSON)
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

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

	return &parsed, nil
}

// Authenticate with username/password api call
// This function sends a request to the MaaS360 authentication API to get an auth token.
func authentication(billingID string, appID string, accessKey string, username string, password string) (*AuthResponse, error) {
	if len(billingID) == 0 {
		return nil, fmt.Errorf("billing ID cannot be empty")
	}
	serviceURL, err := GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}

	auth := maaS360AdminAuth{
		BillingID:  billingID,
		AppID:      appID,
		PlatformID: constants.Platform,
		AppVersion: constants.Version,
		AccessKey:  accessKey,
		Username:   username,
		Password:   password,
	}

	return doAuthRequest(serviceURL, billingID, "/auth-apis/auth/2.0/authenticate", auth, nil)
}

// Refresh token authentication api call
// This function sends a request to the MaaS360 authentication API to refresh an existing auth token.
func refreshToken(billingID string, appID string, accessKey string, username string, refreshToken string) (*AuthResponse, error) {
	if len(billingID) == 0 {
		return nil, fmt.Errorf("billing ID cannot be empty")
	}
	serviceURL, err := GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}

	auth := maaS360AdminAuth{
		BillingID:    billingID,
		AppID:        appID,
		PlatformID:   constants.Platform,
		AppVersion:   constants.Version,
		AccessKey:    accessKey,
		Username:     username,
		RefreshToken: refreshToken,
	}

	return doAuthRequest(serviceURL, billingID, "/auth-apis/auth/2.0/refreshToken", auth, nil)
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

// GetToken retrieves an authentication token from MaaS360.
// It can authenticate using a username/password or refresh an existing token.
// It returns the auth token, refresh token, and any error encountered.
func GetToken(billingID, appID, accessKey, username, password, refresh string) (string, string, error) {
	var response *AuthResponse
	var err error

	if password != "" {
		response, err = authentication(billingID, appID, accessKey, username, password)
		if err != nil {
			log.Printf("Authentication failed: %v", err)
		} else {
			printTokens(*response)
		}
	} else if refresh != "" {
		response, err = refreshToken(billingID, appID, accessKey, username, refresh)
		if err != nil {
			log.Printf("Refresh token failed: %v", err)
		} else {
			printTokens(*response)
		}
	} else {
		return "", "", fmt.Errorf("either password or refresh token must be provided")
	}

	return response.Wrapper.AuthToken, response.Wrapper.RefreshToken, nil
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

// printTokens prints the authentication tokens from the response
// It is used to display the auth token and refresh token if available.
func printTokens(resp AuthResponse) {
	fmt.Printf("Auth Token: %s\n", resp.Wrapper.AuthToken)
	if resp.Wrapper.RefreshToken != "" {
		fmt.Printf("Refresh Token: %s\n", resp.Wrapper.RefreshToken)
	}
}
