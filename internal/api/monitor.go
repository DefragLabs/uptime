package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/defraglabs/uptime/internal/utils"
	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
)

// AddMonitoringURLHandler api lets an user add an healthcheck url.
func AddMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	decoder := json.NewDecoder(r.Body)
	var monitorURLForm forms.MonitorURLForm
	monitorURLForm.UserID = user.ID
	err := decoder.Decode(&monitorURLForm)
	if err != nil {
		writeErrorResponse(w, "Invalid input format")

		log.Info("Invalid input format for forgot password")
		return
	}

	validationMessage := monitorURLForm.Validate()
	if validationMessage != "" {
		writeErrorResponse(w, validationMessage)

		log.Info("Validation failed while adding monitoring URL.")
		return
	}

	datastore := db.New()
	var monitoringURL db.MonitorURL
	monitoringURL = datastore.GetMonitoringURLByUserIDAndURL(user.ID, monitorURLForm.Protocol, monitorURLForm.URL)

	if monitoringURL.ID != "" {
		writeErrorResponse(w, "URL already exists.")

		return
	}

	objectID := db.GenerateObjectID()
	monitorURLForm.ID = objectID.Hex()

	monitoringURL = datastore.AddMonitoringURL(monitorURLForm)
	initialPingMonitorURL(monitoringURL, datastore)

	log.Info(fmt.Sprintf("Added monitoring url %s", monitorURLForm.URL))

	// This is done to retreive the ping results of the added monitor URL.
	monitoringURL = datastore.GetMonitoringURLByUserID(user.ID, monitoringURL.ID)

	responseData := structs.Map(monitoringURL)
	writeSuccessStructResponse(w, responseData, http.StatusCreated)
}

func initialPingMonitorURL(monitorURL db.MonitorURL, datastore *db.Datastore) {
	start := time.Now()
	url := fmt.Sprintf("%s://%s", monitorURL.Protocol, monitorURL.URL)

	resp, err := http.Get(url)
	duration := time.Since(start)
	if err != nil {
		// Don't fail like this.
		log.Warn("API ping failed")
	}

	t := time.Now()
	timeStamp := t.Format(time.UnixDate)

	serviceStatus := utils.GetServiceStatus(resp.StatusCode)
	datastore.AddMonitorDetail(monitorURL, resp.Status, serviceStatus, timeStamp, duration.String())
}

// GetMonitoringURLsHandler api returns the monitoring urls configured
func GetMonitoringURLsHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	monitoringURLS := datastore.GetMonitoringURLSByUserID(user.ID)

	data := make(map[string]interface{})
	data["monitoringURLs"] = monitoringURLS

	writeSuccessStructResponse(w, data, http.StatusOK)
}

// GetMonitoringURLHandler gets an individual monitoringURL
func GetMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	vars := mux.Vars(r)
	monitoringURLID := vars["monitoringURLID"]

	if monitoringURLID == "" {
		writeErrorResponse(w, "Monitoring URL id is required")

		return
	}

	datastore := db.New()
	monitoringURL := datastore.GetMonitoringURLByUserID(user.ID, monitoringURLID)
	if monitoringURL.ID == "" {
		writeErrorResponse(w, "Monitoring url not found")

		return
	}

	responseData := structs.Map(monitoringURL)
	writeSuccessStructResponse(w, responseData, http.StatusOK)
}

// UpdateMonitoringURLHandler api lets the user update the details.
// Can update the following
//   - Protocol
//   - Frequency
//   - Unit
func UpdateMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	vars := mux.Vars(r)
	monitoringURLID := vars["monitoringURLID"]

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	monitoringURL := datastore.GetMonitoringURLByUserID(user.ID, monitoringURLID)
	if monitoringURL.ID == "" {
		writeErrorResponse(w, "Monitoring url not found")

		return
	}

	decoder := json.NewDecoder(r.Body)
	var monitorURLForm forms.MonitorURLForm
	monitorURLForm.UserID = user.ID
	err := decoder.Decode(&monitorURLForm)
	if err != nil {
		writeErrorResponse(w, "Invalid input format")

		log.Info("Invalid input format for forgot password")
		return
	}

	monitorURLForm.ID = monitoringURLID
	monitorURLForm.URL = monitoringURL.URL

	monitoringURL = datastore.UpdateMonitoringURLByUserID(user.ID, monitoringURLID, monitorURLForm)

	responseData := structs.Map(monitoringURL)
	writeSuccessStructResponse(w, responseData, http.StatusOK)
}

// DeleteMonitoringURLHandler api can be used to delete a monitor url
func DeleteMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	vars := mux.Vars(r)
	monitoringURLID := vars["monitoringURLID"]

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	monitoringURL := datastore.GetMonitoringURLByUserID(user.ID, monitoringURLID)
	if monitoringURL.ID == "" {
		writeErrorResponse(w, "Monitoring url not found")

		return
	}

	datastore.DeleteMonitoringURL(user.ID, monitoringURLID)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	log.Info("Monitoring url removed successfully")
}
