package forms

// MonitorURLForm struct represents a row in db.
type MonitorURLForm struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	Protocol  string `bson:"protocol" json:"protocol"`
	URL       string `bson:"url" json:"url"`
	Frequency int32  `bson:"frequency" json:"frequency"`
	Unit      string `bson:"unit" json:"unit"`
}

// Validate monitor url form input
func (monitorURLForm MonitorURLForm) Validate() string {
	if monitorURLForm.Protocol == "" {
		return "Protocol is required"
	} else if monitorURLForm.URL == "" {
		return "URL is required"
	} else if monitorURLForm.Frequency == 0 {
		return "Frequency is required"
	} else if monitorURLForm.Unit == "" {
		return "Unit is required"
	}
	return ""
}
