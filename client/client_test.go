package client

import (
	auth_api "maas360api/auth"
	"maas360api/internal/constants"
	"testing"
)

// TestAuthenticate verifies that the client can be created and authenticated without errors
func TestAuthenticate(t *testing.T) {
	authCredentials := auth_api.MaaS360AdminAuth{
		BillingID:  "123456",
		AppID:      "testApp",
		PlatformID: constants.Platform,
		AppVersion: constants.Version,
		AccessKey:  "testKey",
		Username:   "testUser",
		Password:   "MaaS360!",
		// RefreshToken: "testRefresh", // Optional
	}
	client, err := Authenticate(authCredentials)
	if err != nil {
		t.Fatalf("Expected authentication to succeed, got error: %v", err)
	}
	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}
	if client.BillingID != "123456" {
		t.Errorf("Expected BillingID to be '123456', got '%s'", client.BillingID)
	}
	if client.AppID != "testApp" {
		t.Errorf("Expected AppID to be 'testApp', got '%s'", client.AppID)
	}
	if client.Username != "testUser" {
		t.Errorf("Expected Username to be 'testUser', got '%s'", client.Username)
	}
}

// TestGetBasicAuth verifies basic auth generation
func TestGetBasicAuth(t *testing.T) {
	authCredentials := auth_api.MaaS360AdminAuth{
		BillingID:  "123456",
		AppID:      "testApp",
		PlatformID: constants.Platform,
		AppVersion: constants.Version,
		AccessKey:  "testKey",
		Username:   "testUser",
		Password:   "testPass",
	}
	client, err := Authenticate(authCredentials)
	if err != nil {
		t.Fatalf("Expected authentication to succeed, got error: %v", err)
	}
	basicAuth := client.GetBasicauth()
	expected := "Basic dGVzdFVzZXI6dGVzdFBhc3M="
	if basicAuth != expected {
		t.Errorf("Expected basic auth to be '%s', got '%s'", expected, basicAuth)
	}
}

// TestEmptyBasicAuth verifies empty basic auth handling
func TestEmptyBasicAuth(t *testing.T) {
	authCredentials := auth_api.MaaS360AdminAuth{
		BillingID:  "123456",
		AppID:      "testApp",
		PlatformID: constants.Platform,
		AppVersion: constants.Version,
		AccessKey:  "testKey",
		Username:   "",
		Password:   "",
	}
	client, err := Authenticate(authCredentials)
	if err != nil {
		t.Fatalf("Expected authentication to succeed, got error: %v", err)
	}
	basicAuth := client.GetBasicauth()
	if basicAuth != "" {
		t.Errorf("Expected empty basic auth for empty credentials, got '%s'", basicAuth)
	}
}
