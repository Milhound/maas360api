package devices_api

import (
	"fmt"
	auth_api "maas360api/auth"
)

// UpdateOS schedules an OS update for a specific device in MaaS360.
func UpdateOS(billingID string, deviceID string, osVersion string, targetLocalTime string, maasToken string) error {
	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		fmt.Printf("Error getting serviceURL: %v\n", err)
	}

	detailsURL := serviceURL + "/emc/?#"
	additionalParams := map[string]string{
		"productVersion":     osVersion,
		"osUpdateActionType": "OS Enforcement",
		"targetLocalTime":    targetLocalTime,
		"detailsURL":         detailsURL,
	}

	return PerformDeviceAction(billingID, deviceID, "MDM_SCHEDULE_OS_UPDATE", additionalParams, maasToken)
}
