package devices

import (
	httputil "maas360api/internal/http"
)

type DeviceAttribute struct {
	AttributeKey   string `json:"key"`
	AttributeType  string `json:"type"`
	AttributeValue any    `json:"value"`
}

type DeviceAttributesWrapper struct {
	DeviceAttributes []DeviceAttribute `json:"deviceAttribute"`
}

type DeviceAttributesResponse struct {
	DeviceID         string                  `json:"maas360DeviceID"`
	AttributeWrapper DeviceAttributesWrapper `json:"deviceAttributes"`
}

var client = httputil.GetSharedClient()
