package db

import "github.com/mongodb/mongo-go-driver/bson/objectid"

// MonitorURL struct represents a row in db.
type MonitorURL struct {
	ID        objectid.ObjectID `bson:"-,Skip" json:"id,omitempty"`
	Protocol  string            `bson:"protocol" json:"protocol"`
	URL       string            `bson:"url" json:"url"`
	frequency int32             `bson:"frequency" json:"frequency"`
	unit      string            `bson:"unit" json:"unit"`
}
