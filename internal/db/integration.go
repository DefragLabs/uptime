package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
)

const (
	// EmailIntegration represents email integration.
	EmailIntegration = "email"

	// SlackIntegration represents slack integration.
	SlackIntegration = "slack"

	// PagerDutyIntegration represents pagerduty integration.
	PagerDutyIntegration = "pagerduty"
)

// Integration struct represents a row in db.
type Integration struct {
	ID     string `bson:"_id" json:"id,omitempty" structs:"id"`
	UserID string `bson:"userID" json:"userID" structs:"userID"`

	// Type of the integration.
	// Supported integrations are Slack, Email, PagerDuty.
	Type       string `bson:"type" json:"type" structs:"type"`
	Email      string `bson:"email" json:"email" structs:"email"`
	WebhookURL string `bson:"webhookURL" json:"webhookURL" structs:"webhookURL"`

	// PDRoutingKey is the routing key generated from PD integration.
	PDRoutingKey string `bson:"pdRoutingKey" json:"pdRoutingKey" structs:"pdRoutingKey"`

	// PDAction is the PD event action. Possible values are trigger, acknowledge and resolve.
	PDAction string `bson:"pdAction" json:"pdAction" structs:"pdAction"`

	// PDSeverity can be one of info, warning, error, critical, or unknown.
	PDSeverity string `bson:"pdSeverity" json:"pdSeverity" structs:"pdSeverity"`
}

type slackNotificationMsg struct {
	URL     string
	Message string
}

// Send decides which integration to send notification and sends it.
func (integration *Integration) Send(monitorURL MonitorURL, serviceStatus string) {
	log.Info("Sending alert", integration.Type)

	var err error
	if integration.Type == "slack" {
		log.Info("Try and send slack notification")
		err = integration.SendSlackNotification(monitorURL, serviceStatus)
	} else if integration.Type == "pagerduty" {
		err = integration.SendPagerDutyEvent(monitorURL, serviceStatus)
	} else {
		return
	}

	if err != nil {
		log.Infof("Unable to send integration [%s]", integration.Type)
		return
	}

	log.Infof("Integration %s sent for site %s", integration.Type, monitorURL.URL)
}

// SendPagerDutyEvent sends an event v2 to pagerduty
func (integration *Integration) SendPagerDutyEvent(monitorURL MonitorURL, serviceStatus string) error {
	if integration.PDRoutingKey == "" {
		log.Infof("Invalid integration. PDRoutingKey not found.")

		return errors.New("invalid integration. PDRoutingKey not found")
	}

	timestamp := time.Now().String()
	payload := pagerduty.V2Payload{
		Summary:   fmt.Sprintf("Site %s %s", monitorURL.URL, serviceStatus),
		Source:    "Uptime",
		Severity:  integration.PDSeverity,
		Timestamp: timestamp,
	}
	event := pagerduty.V2Event{
		Payload:    &payload,
		RoutingKey: integration.PDRoutingKey,
		Action:     integration.PDAction,
	}

	_, err := pagerduty.ManageEvent(event)

	if err != nil {
		log.Infof("Pagerduty event send failed for integration %s", integration.ID)

		return errors.New("pagerduty event send failed")
	}

	return nil
}

// SendSlackNotification sends a notification to slack using slack webhooks.
func (integration *Integration) SendSlackNotification(monitorURL MonitorURL, serviceStatus string) error {
	log.Info("Sending slack notification.")

	if integration.WebhookURL == "" {
		log.Infof("Invalid integration. Webhook url not found for integration %s", integration.ID)

		return errors.New("invalid integration. webhook url not found")
	}
	msg := slackNotificationMsg{
		URL:     monitorURL.URL,
		Message: fmt.Sprintf("Site %s %s", monitorURL.URL, serviceStatus),
	}
	byte, _ := json.Marshal(msg)
	_, err := http.Post(integration.WebhookURL, "application/json", bytes.NewBuffer(byte))

	if err != nil {
		return errors.New("slack notification send failed")
	}

	return nil
}
