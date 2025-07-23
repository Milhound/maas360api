package client

import (
	"testing"
)

// TestNewClient verifies that the client can be created without errors
func TestNewClient(t *testing.T) {
	client := NewClient("123456", "testApp", "testKey", "testUser", "testPass", "testRefresh")
	
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
	client := NewClient("123456", "testApp", "testKey", "testUser", "testPass", "testRefresh")
	
	basicAuth := client.GetBasicauth()
	
	// Basic auth should be base64 encoded "testUser:testPass" with "Basic " prefix
	expected := "Basic dGVzdFVzZXI6dGVzdFBhc3M="
	if basicAuth != expected {
		t.Errorf("Expected basic auth to be '%s', got '%s'", expected, basicAuth)
	}
}

// TestEmptyBasicAuth verifies empty basic auth handling
func TestEmptyBasicAuth(t *testing.T) {
	client := NewClient("123456", "testApp", "testKey", "", "", "testRefresh")
	
	basicAuth := client.GetBasicauth()
	
	if basicAuth != "" {
		t.Errorf("Expected empty basic auth for empty credentials, got '%s'", basicAuth)
	}
}