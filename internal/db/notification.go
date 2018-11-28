package db

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Integration struct represents a row in db.
type Integration struct {
	ID         objectid.ObjectID `bson:"-,Skip" json:"id,omitempty"`
	Type       string            `bson:"type" json:"type"`
	EmailID    string            `bson:"emailID" json:"emailID"`
	WebhookURL string            `bson:"webhookURL" json:"webhookURL"`
}
