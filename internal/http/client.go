package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"maas360api/internal/constants"
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

// RequestOptions contains options for making HTTP requests
type RequestOptions struct {
	Method      string
	URL         string
	Body        io.Reader
	ContentType string
	MaaSToken   string
	Context     context.Context
}

// DoMaaSRequest performs a standard MaaS360 API request with proper headers
func DoMaaSRequest(opts RequestOptions) (*http.Response, error) {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}
	
	req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, opts.Body)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Set common headers
	req.Header.Set(constants.AcceptHeader, constants.ContentTypeJSON)
	req.Header.Set(constants.AuthorizationHeader, fmt.Sprintf(constants.MaaSTokenPrefix, opts.MaaSToken))

	// Set content type if provided
	if opts.ContentType != "" {
		req.Header.Set(constants.ContentTypeHeader, opts.ContentType)
	} else {
		req.Header.Set(constants.ContentTypeHeader, constants.ContentTypeJSON)
	}

	resp, err := GetSharedClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	return resp, nil
}
