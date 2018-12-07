package db

// MonitorURL struct represents a row in db.
type MonitorURL struct {
	ID string `bson:"_id" json:"id,omitempty"`

	UserID string `bson:"userID" json:"userID"`

	// Http protocol (http/https)
	Protocol string `bson:"protocol" json:"protocol"`

	// URL that should be pinged.
	URL string `bson:"url" json:"url"`

	// Frequency in integer
	Frequency int32 `bson:"frequency" json:"frequency"`

	// Time unit (minutes, hours)
	Unit string `bson:"unit" json:"unit"`

	// Status of the service. It can be (UP, DOWN, "")
	Status string `bson:"status" json:"status"`

	Results []MonitorResult `bson:"results,omitempty" json:"results"`
}

// MonitorResult contains the ping result.
type MonitorResult struct {
	// Status code of the response.
	Status string

	// Duration of response
	Duration string

	// Timestamp when the ping was run.
	Time string
}
