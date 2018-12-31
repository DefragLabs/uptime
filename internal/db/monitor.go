package db

// MonitorURL struct represents a row in db.
type MonitorURL struct {
	ID string `bson:"_id" json:"id,omitempty" structs:"id"`

	UserID string `bson:"userID" json:"userID" structs:"userID"`

	// Name of the url
	Name string `bson:"name" json:"name" structs:"name"`

	// Http protocol (http/https)
	Protocol string `bson:"protocol" json:"protocol" structs:"protocol"`

	// URL that should be pinged.
	URL string `bson:"url" json:"url" structs:"url"`

	// Frequency in integer
	Frequency int32 `bson:"frequency" json:"frequency" structs:"frequency"`

	// Time unit (minutes, hours)
	Unit string `bson:"unit" json:"unit" structs:"unit"`

	// Status of the service. It can be (UP, DOWN, "")
	Status string `bson:"status" json:"status" structs:"status"`
}

// MonitorResult contains the ping result.
type MonitorResult struct {
	ID string `bson:"_id" json:"id,omitempty" structs:"id"`

	// MonitorURLID stores the relationship to MonitorURL
	MonitorURLID string `bson:"monitorURLID" json:"monitorURLID" structs:"monitorURLID"`

	// Status UP/Down of the response.
	Status string `bson:"status" json:"status" structs:"status"`

	// Status code of the response.
	StatusDescription string `bson:"statusDescription" json:"statusDescription" structs:"statusDescription"`

	// Response time
	ResponseTime float64 `bson:"responseTime" json:"responseTime" structs:"responseTime"`

	// Timestamp when the ping was run.
	Time string `bson:"time" json:"time" structs:"time"`
}
