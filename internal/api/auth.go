package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/defraglabs/uptime/internal/utils"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

// RegisterHandler registers the user.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var userRegisterForm forms.UserRegisterForm
	err := decoder.Decode(&userRegisterForm)

	errorMsg := ""
	if err != nil {
		writeErrorResponse(w, "Invalid input format")

		return
	}

	validationMessage := userRegisterForm.Validate()
	if validationMessage != "" {
		writeErrorResponse(w, errorMsg)

		return
	}

	datastore := db.New()
	user := datastore.GetUserByEmail(userRegisterForm.Email)
	if user.ID != "" {
		errorMsg = fmt.Sprintf("Email %s already registered", user.Email)
		writeErrorResponse(w, errorMsg)

		return
	}

	companyUser := datastore.GetUserByComapnyName(userRegisterForm.CompanyName)

	if companyUser.ID != "" {
		writeErrorResponse(w, fmt.Sprintf("Company %s already exists", companyUser.CompanyName))

		return
	}

	newUser := db.RegisterUser(userRegisterForm)

	objectID := db.GenerateObjectID()
	newUser.ID = objectID.Hex()

	datastore.CreateUser(newUser)

	log.Infof("Registration successful with email %s", newUser.Email)
	w.WriteHeader(http.StatusCreated)
}

// LoginHandler validates the password & returns the JWT token.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)

	var userLoginForm forms.UserLoginForm
	err := decoder.Decode(&userLoginForm)

	error := false
	errorMsg := ""

	if err != nil {
		error = true
		errorMsg = "Invalid input format"
	}

	datastore := db.New()

	user := datastore.GetUserByEmail(userLoginForm.Email)
	fmt.Print(user, user.ID)
	if user.ID == "" {
		error = true
		errorMsg = "User not found"
	}

	validationMessage := userLoginForm.Validate()
	if validationMessage != "" {
		error = true
		errorMsg = validationMessage
	}

	if error {
		errorVal := make(map[string]string)
		errorVal["message"] = errorMsg
		response := Response{
			Success: false,
			Data:    nil,
			Error:   errorVal,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		log.Info("Login failed")
		return
	}

	jwt := db.GetJWT(user, userLoginForm.Password)
	data := make(map[string]string)
	data["token"] = jwt

	response := Response{
		Success: true,
		Data:    data,
		Error:   nil,
	}
	w.WriteHeader(http.StatusOK)

	log.Info("Login successful")
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler revokes the jwt token.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	_, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	if db.RevokeToken(authToken) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		writeErrorResponse(w, "Unable to revoke the token")
	}
}

// ForgotPasswordHandler sends forgot password email.
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var forgotPasswordForm forms.ForgotPasswordForm
	err := decoder.Decode(&forgotPasswordForm)

	if err != nil {
		writeErrorResponse(w, "Invalid input format")

		log.Info("Invalid input format for forgot password")
		return
	}

	validationMessage := forgotPasswordForm.Validate()
	if validationMessage != "" {
		writeErrorResponse(w, validationMessage)

		log.Info("Unable to process forgot password.")
		return
	}

	datastore := db.New()
	user := datastore.GetUserByEmail(forgotPasswordForm.Email)

	if user.ID == "" {
		writeErrorResponse(w, "User not found")

		log.Infof("User with email %s not found", forgotPasswordForm.Email)
		return
	}

	toEmail := user.Email

	code := uuid.Must(uuid.NewV4())
	hexCode := hex.EncodeToString(code.Bytes())

	resetPassword := db.ResetPassword{UserID: user.ID, Code: hexCode}

	datastore.AddResetPassword(resetPassword)

	baseURL, _ := url.Parse(fmt.Sprintf("http://%s", r.Host))
	baseURL.Path = path.Join(baseURL.Path, os.Getenv("FORGOT_PASSWORD_LINK"), user.ID, hexCode)

	forgotPasswordLink := baseURL.String()

	sub := "Forgot password"
	msg := fmt.Sprintf(
		"Hi,"+
			"Click this <a href=%s></a>"+
			"\r\n", forgotPasswordLink,
	)

	utils.SendMail(sub, msg, toEmail)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	log.Infof("Successfully sent forgot password mail to %s", toEmail)
}

// ResetPasswordHandler password reset handler
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var resetPasswordForm forms.ResetPasswordForm
	err := decoder.Decode(&resetPasswordForm)
	if err != nil {
		panic(err)
	}

	success := db.PasswordReset(resetPasswordForm.UID, resetPasswordForm.Code, resetPasswordForm.NewPassword)
	datastore := db.New()

	if success == true {
		user := datastore.GetUserByID(resetPasswordForm.UID)

		sub := "Forgot password"
		msg := fmt.Sprintf(
			"Hi," +
				"Your password has been reset successfully.",
		)

		utils.SendMail(sub, msg, user.Email)
	} else {
		writeErrorResponse(w, "Password cannot be reset. Please try again.")

		log.Info("Unable to reset password.")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	log.Info("Password reset success.")
}
