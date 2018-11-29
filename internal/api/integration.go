package api

import (
	"encoding/json"
	"net/http"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// AddIntegrationHandler can be used to add a new integration.
func AddIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	decoder := json.NewDecoder(r.Body)
	var integrationForm forms.IntegrationForm
	integrationForm.UserID = user.ID
	err := decoder.Decode(&integrationForm)
	if err != nil {
		writeErrorResponse(w, "Invalid input format")

		log.Info("Invalid input format for forgot password")
		return
	}

	validationMessage := integrationForm.Validate()
	if validationMessage != "" {
		writeErrorResponse(w, validationMessage)

		log.Info("Validation failed while adding integration.")
		return
	}

	objectID := db.GenerateObjectID()
	integrationForm.ID = objectID.Hex()

	datastore := db.New()
	integration := datastore.AddIntegration(integrationForm)

	log.Info("Integration added successfully.")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(integration)
}

// GetIntegrationsHandler gets all integrations by the logged in user.
func GetIntegrationsHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	integrations := datastore.GetIntegrationsByUserID(user.ID)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(integrations)
}

// GetIntegrationHandler gets a specific integration
func GetIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	integrationID := vars["integrationID"]

	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	integration := datastore.GetIntegrationByUserID(user.ID, integrationID)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(integration)
}

// DeleteIntegrationHandler removes a configured integration
func DeleteIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	integrationID := vars["integrationID"]

	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	datastore.DeleteIntegration(user.ID, integrationID)

	log.Info("Integration removed successfully.")
}
