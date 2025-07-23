# 📦 maas360api (Go Library)

## Overview

`maas360api` is a Go client library for interacting with the IBM MaaS360 Web Services API. It supports authentication and encapsulates API operations for managing applications and devices within a MaaS360 environment.

## ✨ Features

- Secure token-based authentication with MaaS360
- Device API integration (search, retrieve, manage)
- Application API integration (list, upload, remove)
- Modular, extensible, idiomatic Go design

## 📥 Installation

```bash
go get github.com/Milhound/maas360api
````

## 🔐 Authentication

```go
package main

import (
    "log"

    "maas360api/auth"
    "maas360api/client"
    "maas360api/internal/constants"
)

func main() {
    authCredentials := auth_api.MaaS360AdminAuth{
        BillingID:  "<YOUR_BILLING_ID>",
        AppID:      "<YOUR_APP_ID>",
        PlatformID: constants.Platform,
        AppVersion: constants.Version,
        AccessKey:  "<YOUR_ACCESS_KEY>",
        Username:   "<YOUR_USERNAME>",
        Password:   "<YOUR_PASSWORD>",
        // RefreshToken: "<YOUR_REFRESH_TOKEN>", // Optional, can be empty if not using refresh token
    }

    MaaS360, err := MaaS360api.Authenticate(authCredentials)
    if err != nil {
        log.Fatalf("Error authenticating: %v", err)
    }

    // MaaS360 client is now authenticated and ready for use
}

```

## 🙌 Contributing

Contributions are welcome! Please:

* Open issues for bugs or feature requests
* Fork and submit pull requests for enhancements
* Follow idiomatic Go and clean code principles
