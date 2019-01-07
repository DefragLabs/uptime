package forms

import (
	"fmt"
	"net/http"

	"github.com/defraglabs/uptime/internal/utils"
)

// MonitorURLForm struct represents a row in db.
type MonitorURLForm struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	UserID    string `bson:"userID" json:"-"`
	Name      string `bson:"name" json:"name"`
	Protocol  string `bson:"protocol" json:"protocol"`
	URL       string `bson:"url" json:"url"`
	Frequency int32  `bson:"frequency" json:"frequency"`
	Unit      string `bson:"unit" json:"unit"`
}

func validateURL(url string) bool {
	resp, err := http.Get(url)

	if err != nil {
		return false
	}

	// Status codes should be in 2xx range.
	if resp.StatusCode >= 300 {
		return false
	}

	return true
}

// Validate monitor url form input
func (monitorURLForm MonitorURLForm) Validate() string {
	// TODO: Uncomment this for the next release.
	// if monitorURLForm.Name == "" {
	// 	return "Name is required"
	// } else
	if monitorURLForm.Protocol == "" {
		return "Protocol is required"
	} else if monitorURLForm.URL == "" {
		return "URL is required"
	} else if monitorURLForm.Frequency == 0 {
		return "Frequency is required"
	} else if monitorURLForm.Unit == "" {
		return "Unit is required"
	}

	// Validate if the provided frequency and units are valid.
	if val, ok := utils.MonitoringConfig[monitorURLForm.Unit]; ok {
		if !utils.FrequencyInMonitoringConfig(monitorURLForm.Frequency, val) {
			return "Invalid frequency"
		}
	} else {
		return "Invalid unit"
	}

	url := fmt.Sprintf("%s://%s", monitorURLForm.Protocol, monitorURLForm.URL)
	if !validateURL(url) {
		return "Make sure you've provided the correct url & protocol."
	}
	return ""
}
