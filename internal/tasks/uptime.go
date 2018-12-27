package tasks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/utils"
	log "github.com/sirupsen/logrus"
)

func pingURL(t time.Time) {
	datastore := db.New()
	monitoringURLS := datastore.GetMonitoringURLS()

	currentTime := time.Now()
	if int32(currentTime.Minute()*60+currentTime.Second())%(5*60) == 0 {
		log.Infof("Start pinging urls. Total urls %d", len(monitoringURLS))
	}

	for _, monitorURL := range monitoringURLS {
		url := fmt.Sprintf("%s://%s", monitorURL.Protocol, monitorURL.URL)
		start := currentTime

		// Validate if the provided frequency and units are valid.
		if val, ok := utils.MonitoringConfig[monitorURL.Unit]; ok {
			if !utils.FrequencyInMonitoringConfig(monitorURL.Frequency, val) {
				log.Infof("Invalid frequency found for url %s", monitorURL.URL)
			}
		} else {
			log.Infof("Invalid unit found for url %s", monitorURL.URL)
			continue
		}

		if monitorURL.Unit == utils.MINUTE && int32(currentTime.Minute()*60+currentTime.Second())%(monitorURL.Frequency*60) != 0 {
			continue
		} else if monitorURL.Unit == utils.SECOND && int32(currentTime.Second())%monitorURL.Frequency != 0 {
			continue
		}

		resp, err := http.Get(url)
		responseTime := time.Since(start)
		if err != nil {
			// Don't fail like this.
			log.Warn("API ping failed")
		}
		timeStamp := t.Format(time.UnixDate)
		fmt.Println(responseTime, url, resp.Status, timeStamp)

		serviceStatus := utils.GetServiceStatus(resp.StatusCode)

		responseTimeInMillSeconds := float64(responseTime.Nanoseconds()) / 1000000
		datastore.AddMonitorDetail(monitorURL, resp.Status, serviceStatus, timeStamp, responseTimeInMillSeconds)

		notify := shouldNotify(monitorURL, serviceStatus)
		if notify {
			sendAlertNotification(monitorURL, serviceStatus)
		}
	}
}

// shouldNotify checks if a notification has to be sent.
func shouldNotify(monitorURL db.MonitorURL, serviceStatus string) bool {
	datastore := db.New()
	results := datastore.GetLastNMonitoringURLStats(monitorURL.ID, 1)
	if len(results) == 1 {
		latestResult := results[0]
		if latestResult.Status != serviceStatus {
			if serviceStatus == utils.StatusUp {
				return true
			} else if serviceStatus == utils.StatusDown {
				return true
			}
		}
	}

	return false
}

// sendAlertNotification sends a notification through all the configured integrations
func sendAlertNotification(monitorURL db.MonitorURL, serviceStatus string) {
	userID := monitorURL.UserID

	datastore := db.New()
	userIntegrations := datastore.GetIntegrationsByUserID(userID)

	for _, integration := range userIntegrations {
		integration.Send(monitorURL, serviceStatus)
	}
}

// StartScheduler runs the scheduler
func StartScheduler() {
	log.Info("Starting scheduler")
	c := make(chan time.Time)
	go func() {
		var frequency = time.Duration(1 * time.Second)
		ticker := time.Tick(frequency)

		for {
			time.Sleep(time.Duration(1 * time.Second))
			c <- <-ticker
		}
	}()

	for {
		select {
		case t := <-c:
			pingURL(t)
		case <-time.After(time.Duration(45 * time.Minute)):
			// This case acts as a timeout.
			log.Info("Task timeout.")
		}
	}
}
