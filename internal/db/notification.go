package db

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Notification struct represents a row in db.
type Notification struct {
	ID         objectid.ObjectID `bson:"_id" json:"id"`
	Type       string            `bson:"type" json:"type"`
	EmailID    string            `bson:"emailID" json:"emailID"`
	webhookURL string            `bson:"webhookURL" json:"webhookURL"`
}
