package http

import (
	"net/http"
	"time"
)

// DefaultTimeout is the default timeout for HTTP requests
const DefaultTimeout = 30 * time.Second

// NewClient creates a new HTTP client with optimized settings for the MaaS360 API
func NewClient() *http.Client {
	return &http.Client{
		Timeout: DefaultTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// GetSharedClient returns a shared HTTP client instance
var sharedClient *http.Client

func GetSharedClient() *http.Client {
	if sharedClient == nil {
		sharedClient = NewClient()
	}
	return sharedClient
}