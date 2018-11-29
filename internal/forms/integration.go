package forms

// IntegrationForm struct is used for input data for integrations.
type IntegrationForm struct {
	ID         string `bson:"_id" json:"id,omitempty"`
	UserID     string `bson:"userID" json:"-"`
	Type       string `bson:"type" json:"type"`
	EmailID    string `bson:"emailID" json:"emailID"`
	WebhookURL string `bson:"webhookURL" json:"webhookURL"`
}

// Validate integration form
func (integrationForm IntegrationForm) Validate() string {
	if integrationForm.Type == "" {
		return "integration type is required"
	} else if integrationForm.Type == "slack" && integrationForm.WebhookURL == "" {
		return "webhookURL is required for slack integration"
	} else if integrationForm.Type == "email" && integrationForm.EmailID == "" {
		return "email id is required for mail integration"
	}

	return ""
}
