package db

// Integration struct represents a row in db.
type Integration struct {
	ID         string `bson:"_id" json:"id,omitempty" structs:"id"`
	UserID     string `bson:"userID" json:"userID" structs:"userID"`
	Type       string `bson:"type" json:"type" structs:"type"`
	Email      string `bson:"email" json:"email" structs:"email"`
	WebhookURL string `bson:"webhookURL" json:"webhookURL" structs:"webhookURL"`
}
