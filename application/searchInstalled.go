package application_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	auth_api "maas360api/auth"
)

type InstalledApp struct {
	AppID         string `json:"appID"`
	AppName       string `json:"appName"`
	DeviceCount   int    `json:"deviceCount"`
	MajorVersions int    `json:"majorVersions"` // <-- change from string to int
	Platform      string `json:"platform"`
}

type InstalledApps struct {
	App        []InstalledApp `json:"app"`
	Count      int            `json:"count"`
	PageSize   int            `json:"pageSize"`
	PageNumber int            `json:"pageNumber"`
}

type InstalledAppsResponse struct {
	InstalledApps InstalledApps `json:"installedApps"`
}

// SearchInstalledApps retrieves installed applications based on the provided filters.
func SearchInstalledApps(billingID, token string, filters map[string]string) ([]InstalledApp, error) {
	// Search parameters: All are optional
	// partialAppName - Partial or full App Name string that needs to be searched for
	// appID - Full AppID that needs to be searched for
	// platform - Supported values: [iOS, Android, BlackBerry]
	// pageSize - Limit number of applications returned at one time. Allowed page sizes: 25, 50, 100, 200, 250. Default value: 25.
	// pageNumber - Results specific to a particular page. Default is first page

	// Validate required fields
	if len(billingID) == 0 || len(token) == 0 {
		return nil, fmt.Errorf("billing ID and token cannot be empty")
	}

	// Get the MaaS360 instance URL
	instance, err := auth_api.GetInstance(billingID)
	if err != nil {
		return nil, err
	}
	if filters == nil {
		searchURL := fmt.Sprintf("%s/application-apis/installedApps/1.0/search/%s?", instance, billingID)
		return doSearchRequest(searchURL, token)
	} else if len(filters) == 0 {
		return nil, fmt.Errorf("search parameters cannot be empty")
	} else {
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
		searchURL := fmt.Sprintf("%s/application-apis/installedApps/1.0/search/%s?", instance, billingID) + searchFilters.Encode()
		return doSearchRequest(searchURL, token)
	}
}

// doSearchRequest sends a request to the MaaS360 API to search for installed applications.
// It constructs the request, sends it, and processes the response.
func doSearchRequest(url string, token string) ([]InstalledApp, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("MaaS token=\"%s\"", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	var InstalledAppsResponse InstalledAppsResponse
	if err := json.Unmarshal(body, &InstalledAppsResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response JSON: %v; body: %s", err, string(body))
	}
	if len(InstalledAppsResponse.InstalledApps.App) == 0 {
		return nil, fmt.Errorf("no apps found")
	}
	return InstalledAppsResponse.InstalledApps.App, nil
}

// PrintAllSoftwareInstalled retrieves and prints all installed software for a given billing ID.
// It uses SearchInstalledApps to get the list of installed applications and formats the output.
func PrintAllSoftwareInstalled(billingID string, token string, filters map[string]string) {
	apps, err := SearchInstalledApps(billingID, token, filters)
	if err != nil {
		log.Fatalf("Error searching installed apps: %v", err)
	}
	if len(apps) == 0 {
		log.Fatal("No installed apps found")
	}
	log.Printf("Found %d installed apps", len(apps))

	for _, app := range apps {
		fmt.Printf("App Name: %s, App ID: %s, Platform: %s, Device Count: %d\n", app.AppName, app.AppID, app.Platform, app.DeviceCount)
	}
}
