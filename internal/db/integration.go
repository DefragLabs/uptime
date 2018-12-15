package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
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

// SendPagerDutyEvent sends an event v2 to pagerduty
func (integration *Integration) SendPagerDutyEvent(timestamp string) error {
	if integration.PDRoutingKey == "" {
		log.Infof("Invalid integration. PDRoutingKey not found.")

		return errors.New("invalid integration. PDRoutingKey not found")
	}

	payload := pagerduty.V2Payload{
		Summary:   "Site down",
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
func (integration *Integration) SendSlackNotification(monitorURL MonitorURL) error {
	if integration.WebhookURL == "" {
		log.Infof("Invalid integration. Webhook url not found for integration %s", integration.ID)

		return errors.New("invalid integration. webhook url not found")
	}
	msg := slackNotificationMsg{
		URL:     monitorURL.URL,
		Message: "Site down",
	}
	byte, _ := json.Marshal(msg)
	_, err := http.Post(integration.WebhookURL, "application/json", bytes.NewBuffer(byte))

	if err != nil {
		return errors.New("slack notification send failed")
	}

	return nil
}
