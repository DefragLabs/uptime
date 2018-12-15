package db

import (
	"errors"

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
	}

	return nil
}
