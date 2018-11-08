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

	error := false
	errorMsg := ""
	if err != nil {
		error = true
		errorMsg = "Invalid input format"
	}

	datastore := db.New()
	user := datastore.GetUserByEmail(userRegisterForm.Email)
	if user.ID != "" {
		error = true
		errorMsg = "Email already registered"
	}

	validationMessage := userRegisterForm.Validate()
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
		log.Info("Registration failed")
		return
	}

	newUser := db.RegisterUser(userRegisterForm)

	objectID := db.GenerateObjectID()
	newUser.ID = objectID.Hex()

	datastore.CreateUser(newUser)

	log.Info("Registration successful")
	w.WriteHeader(http.StatusCreated)
}

// LoginHandler validates the password & returns the JWT token.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)

	var userLoginForm forms.UserLoginForm
	err := decoder.Decode(&userLoginForm)
	if err != nil {
		panic(err)
	}

	datastore := db.New()

	user := datastore.GetUserByEmail(userLoginForm.Email)
	jwt := db.GetJWT(user, userLoginForm.Password)
	data := make(map[string]string)
	data["token"] = jwt

	response := Response{
		Success: true,
		Data:    data,
		Error:   nil,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ForgotPasswordHandler sends forgot password email.
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	decoder := json.NewDecoder(r.Body)
	var forgotPasswordForm forms.ForgotPasswordForm
	err := decoder.Decode(&forgotPasswordForm)
	if err != nil {
		panic(err)
	}

	datastore := db.New()
	user := datastore.GetUserByEmail(forgotPasswordForm.Email)
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
		errorVal := make(map[string]string)
		errorVal["message"] = "Password cannot be reset. Please try again."
		response := Response{
			Success: false,
			Data:    nil,
			Error:   errorVal,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}
	w.WriteHeader(http.StatusOK)
}
