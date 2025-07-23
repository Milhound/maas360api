package application

import (
	httputil "maas360api/internal/http"
)

var client = httputil.GetSharedClient()
