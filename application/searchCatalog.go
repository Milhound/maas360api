package application_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	auth_api "maas360api/auth"
	"maas360api/internal/constants"
)

type CatalogApp struct {
	AppIconUrl       string `json:"appIconURL"`
	AppName          string `json:"appName"`
	AppID            string `json:"appId"`
	EnterpriseRating string `json:"enterpriseRating"`
	FileName         string `json:"fileName"`
	Platform         string `json:"platform"`
	AppType          int    `json:"appType"`
	AppIconFullUrl   string `json:"appIconFullURL"`
	AppFullVersion   string `json:"appFullVersion"`
	AppVersionState  int    `json:"appVersionState"`
	Category         string `json:"category"`
	FileSize         string `json:"fileSize"`
	Status           string `json:"status"`
	DeviceType       int    `json:"deviceType"`
	UploadDate       string `json:"uploadDate"`
	UploadedBy       string `json:"uploadedBy"`
	LastUpdated      string `json:"lastUpdated"`
	InstantUpdate    int    `json:"instantUpdate"`
	LastUpdatedBy    string `json:"lastUpdatedBy"`
	GroupName        string `json:"groupName"`
	GroupId          int    `json:"groupId"`
	SsId             int    `json:"ssId"`
	VppCodes         string `json:"vppCodes"`
}

type CatalogApps struct {
	Count      int          `json:"count"`
	PageSize   int          `json:"pageSize"`
	PageNumber int          `json:"pageNumber"`
	Apps       []CatalogApp `json:"app"`
}

type CatalogAppsResponse struct {
	CatalogApps CatalogApps `json:"apps"`
}

// SearchCatalog retrieves applications from the MaaS360 catalog based on the provided filters.
func SearchCatalog(billingID string, filters map[string]string, maasToken string) ([]CatalogApp, error) {
	// Parameters:
	// pageSize: Limit number of applications returned at one time. Allowed page sizes: 25, 50, 100, 200, 250. Default value: 25.
	// pageNumber: The page number of the results to return. Default value: 1.
	// appName: Partial Application Name string that needs to be searched for.
	// appId: (REQUIRED) Partial or full App ID for the app to be searched.
	// appType: Possible values: 1: iOS Enterprise Application, 2: iOS App Store Application, 3: Android Enterprise Application, 4: Android Market Application, 8: iOS Web-Clip, 10: Mac App Store Application, 11: Mac Enterprise Application
	// enterpriseRating: Possible values: [1, 2, 3, 4, 5] where 1 is the lowest rating and 5 is the highest rating.
	// category: The category of the application. Possible values: [Business, Education, Entertainment, Finance, Health & Fitness, Lifestyle, Medical, Music, Navigation, News, Photography, Productivity, Reference, Social Networking, Sports, Travel & Local, Utilities, Weather].
	// status: Possible values: [Active, Deleted] case insensitive.
	// deviceType: Possible values: 1: Smartphone, 2: Tablet, 3: Smartphone, Tablet
	// instantUpdate: Possible values: 0: Disabled, 1: Enabled.

	// Validate required fields
	if len(billingID) == 0 || len(maasToken) == 0 {
		return nil, fmt.Errorf("billing ID and maasToken cannot be empty")
	}
	// Get the MaaS360 service URL
	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		return nil, err
	}
	if filters == nil {
		return nil, fmt.Errorf("search parameters cannot be nil")
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
		if searchFilters.Get("appId") == "" {
			return nil, fmt.Errorf("appId is a required parameter and cannot be empty")
		}

		searchURL := fmt.Sprintf("%s/application-apis/applications/2.0/search/customer/%s?", serviceURL, billingID) + searchFilters.Encode()
		return doSearchCatalogRequest(searchURL, maasToken)
	}
}

// doSearchCatalogRequest sends a request to the MaaS360 API to search for catalog applications.
// It constructs the request, sends it, and processes the response.
func doSearchCatalogRequest(url string, maasToken string) ([]CatalogApp, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeForm)
	req.Header.Set(constants.AcceptHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.AuthorizationHeader, fmt.Sprintf(constants.MaaSTokenPrefix, maasToken))

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
	var CatalogAppsResponse CatalogAppsResponse
	if err := json.Unmarshal(body, &CatalogAppsResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response JSON: %v; body: %s", err, string(body))
	}
	return CatalogAppsResponse.CatalogApps.Apps, nil
}

// PrintCatalogApps retrieves and prints the catalog applications for a given billing ID.
// It uses SearchCatalog to get the list of applications and formats the output.
func PrintCatalogApps(billingID string, filters map[string]string, maasToken string) {
	apps, err := SearchCatalog(billingID, filters, maasToken)
	if err != nil {
		log.Fatalf("Error searching catalog: %v", err)
	}
	if len(apps) == 0 {
		log.Fatal("No apps found")
	}
	log.Printf("Found %d apps", len(apps))
	for _, app := range apps {
		fmt.Printf("App Name: %s, App ID: %s, Platform: %s, Category: %s, Uploaded By: %s\n", app.AppName, app.AppID, app.Platform, app.Category, app.UploadedBy)
	}
}
