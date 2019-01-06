package api

import (
	"encoding/json"
	"net/http"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/fatih/structs"
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

	responseData := structs.Map(integration)
	writeSuccessStructResponse(w, responseData, http.StatusCreated)
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

	data := make(map[string][]db.Integration)
	// data["integrations"] = integrations
	for _, integration := range integrations {
		if _, ok := data[integration.Type]; ok {
			data[integration.Type] = append(data[integration.Type], integration)
		} else {
			data[integration.Type] = append(data[integration.Type], integration)
		}
	}

	writeSuccessSimpleResponse(w, data, http.StatusOK)
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
	responseData := structs.Map(integration)
	writeSuccessStructResponse(w, responseData, http.StatusOK)
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	log.Info("Integration removed successfully.")
}
