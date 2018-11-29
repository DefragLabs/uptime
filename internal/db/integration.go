package db

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Integration struct represents a row in db.
type Integration struct {
	ID         objectid.ObjectID `bson:"-,Skip" json:"id,omitempty"`
	UserID     string            `bson:"userID" json:"userID"`
	Type       string            `bson:"type" json:"type"`
	Email      string            `bson:"email" json:"email"`
	WebhookURL string            `bson:"webhookURL" json:"webhookURL"`
}
