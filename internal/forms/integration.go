package forms

// IntegrationForm struct is used for input data for integrations.
type IntegrationForm struct {
	ID         string `bson:"_id" json:"id,omitempty"`
	UserID     string `bson:"userID" json:"userID,omitempty"`
	Type       string `bson:"type" json:"type"`
	Email      string `bson:"email" json:"email"`
	WebhookURL string `bson:"webhookURL" json:"webhookURL"`

	// PDRoutingKey is the routing key generated from PD integration.
	PDRoutingKey string `bson:"pdRoutingKey" json:"pdRoutingKey,omitempty"`

	// PDAction is the PD event action. Possible values are trigger, acknowledge and resolve.
	PDAction string `bson:"pdAction" json:"pdAction,omitempty" structs:"pdAction"`

	// PDSeverity can be one of info, warning, error, critical, or unknown.
	PDSeverity string `bson:"pdSeverity" json:"pdSeverity,omitempty"`
}

// Validate integration form
func (integrationForm IntegrationForm) Validate() string {
	if integrationForm.Type == "" {
		return "integration type is required"
	} else if integrationForm.Type == "slack" && integrationForm.WebhookURL == "" {
		return "webhookURL is required for slack integration"
	} else if integrationForm.Type == "email" && integrationForm.Email == "" {
		return "email id is required for mail integration"
	} else if integrationForm.Type == "pagerduty" && (integrationForm.PDRoutingKey == "" || integrationForm.PDAction == "" || integrationForm.PDSeverity == "") {
		return "pagerduty routing key, action & severity are required"
	}

	return ""
}
