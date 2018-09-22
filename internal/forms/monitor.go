package forms

// MonitorURL struct represents a row in db.
type MonitorURL struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	Protocol  string `bson:"protocol" json:"protocol"`
	URL       string `bson:"url" json:"url"`
	Frequency int32  `bson:"frequency" json:"frequency"`
	Unit      string `bson:"unit" json:"unit"`
}
