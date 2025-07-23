package devices

import (
	"fmt"
	auth_api "maas360api/auth"
	"time"
)

// UpdateOS schedules an OS update for a specific device in MaaS360.
func UpdateOS(billingID string, deviceID string, osVersion string, targetLocalTime time.Time, maasToken string) error {
	if billingID == "" || deviceID == "" || osVersion == "" || targetLocalTime.Equal((time.Time{})) || maasToken == "" {
		return fmt.Errorf("billingID, deviceID, osVersion, targetLocalTime, and maasToken must not be empty")
	}

	serviceURL, err := auth_api.GetServiceURL(billingID)
	if err != nil {
		fmt.Printf("Error getting serviceURL: %v\n", err)
	}

	formattedTime := targetLocalTime.Format("2006-01-02T15:04:05")
	detailsURL := serviceURL + "/emc/?#"

	additionalParams := map[string]string{
		"productVersion":     osVersion,
		"osUpdateActionType": "OS Enforcement",
		"targetLocalTime":    formattedTime,
		"detailsURL":         detailsURL,
	}

	return PerformDeviceAction(billingID, deviceID, "MDM_SCHEDULE_OS_UPDATE", additionalParams, maasToken)
}
