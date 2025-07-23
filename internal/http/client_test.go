package http

import (
	"net/http"
	"testing"
	"time"
)

// TestNewClient verifies that a new HTTP client is created with proper settings
func TestNewClient(t *testing.T) {
	client := NewClient()
	
	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}
	
	if client.Timeout != DefaultTimeout {
		t.Errorf("Expected timeout to be %v, got %v", DefaultTimeout, client.Timeout)
	}
	
	// Check transport settings
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected client to have http.Transport")
	}
	
	if transport.MaxIdleConns != 100 {
		t.Errorf("Expected MaxIdleConns to be 100, got %d", transport.MaxIdleConns)
	}
	
	if transport.MaxIdleConnsPerHost != 10 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 10, got %d", transport.MaxIdleConnsPerHost)
	}
	
	if transport.IdleConnTimeout != 90*time.Second {
		t.Errorf("Expected IdleConnTimeout to be 90s, got %v", transport.IdleConnTimeout)
	}
}

// TestGetSharedClient verifies that the shared client is reused
func TestGetSharedClient(t *testing.T) {
	// Reset shared client for test
	sharedClient = nil
	
	client1 := GetSharedClient()
	client2 := GetSharedClient()
	
	if client1 != client2 {
		t.Error("Expected shared client to return the same instance")
	}
	
	if client1 == nil {
		t.Fatal("Expected shared client to be created, got nil")
	}
}