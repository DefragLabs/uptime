package api

import (
	"encoding/json"
	"net/http"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
)

// GetUserDetailHandler returns user details.
func GetUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	responseData := structs.Map(user)
	writeSuccessStructResponse(w, responseData, http.StatusOK)
}

// UpdateUserDetailHandler updates the user details.
func UpdateUserDetailHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	decoder := json.NewDecoder(r.Body)
	var userDetailForm forms.UserDetailForm

	userDetailForm.FirstName = user.FirstName
	userDetailForm.LastName = user.LastName
	userDetailForm.Email = user.Email
	userDetailForm.CompanyName = user.CompanyName
	userDetailForm.PhoneNumber = user.PhoneNumber

	err := decoder.Decode(&userDetailForm)
	if err != nil {
		writeErrorResponse(w, "Invalid input format")

		log.Info("Invalid input format for forgot password")
		return
	}

	userDetailForm.ID = user.ID

	datastore := db.New()
	datastore.UpdateUserDetail(user.ID, userDetailForm)

	// Refresh from database.
	user = datastore.GetUserByID(user.ID)
	responseData := structs.Map(user)
	writeSuccessStructResponse(w, responseData, http.StatusOK)
}
