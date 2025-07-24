package devices

import (
	"fmt"
	"time"
)

// UpdateOS schedules an OS update for a specific device in MaaS360.
func UpdateOS(serviceURL string, billingID string, deviceID string, osVersion string, targetLocalTime time.Time, maasToken string) error {
	if serviceURL == "" || billingID == "" || deviceID == "" || osVersion == "" || targetLocalTime.Equal((time.Time{})) || maasToken == "" {
		return fmt.Errorf("serviceURL, billingID, deviceID, osVersion, targetLocalTime, and maasToken must not be empty")
	}

	formattedTime := targetLocalTime.Format("2006-01-02T15:04:05")
	detailsURL := serviceURL + "/emc/?#"

	additionalParams := map[string]string{
		"productVersion":     osVersion,
		"osUpdateActionType": "OS Enforcement",
		"targetLocalTime":    formattedTime,
		"detailsURL":         detailsURL,
	}

	return PerformDeviceAction(serviceURL, billingID, deviceID, "MDM_SCHEDULE_OS_UPDATE", additionalParams, maasToken)
}
