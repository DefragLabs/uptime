package api

import (
	"net/http"

	"github.com/defraglabs/uptime/internal/utils"

	"github.com/defraglabs/uptime/internal/db"
)

// DashboardStatsHandler returns dashboard stats.
func DashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	monitoringURLCount := datastore.GetMonitoringURLSByUserIDCount(user.ID)
	upMonitoringURLCount := datastore.GetMonitoringURLSByUserIDAndStatus(user.ID, utils.StatusUp)
	downMonitoringURLCount := datastore.GetMonitoringURLSByUserIDAndStatus(user.ID, utils.StatusDown)

	stats := make(map[string]interface{})

	stats["monitoring_urls_count"] = monitoringURLCount
	stats["up_monitoring_urls_count"] = upMonitoringURLCount
	stats["down_monitoring_urls_count"] = downMonitoringURLCount

	writeSuccessStructResponse(w, stats, http.StatusOK)
}
