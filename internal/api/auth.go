package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/gofrs/uuid"

	"github.com/dineshs91/uptime/internal/db"
	"github.com/dineshs91/uptime/internal/forms"
	"github.com/dineshs91/uptime/internal/utils"
)

// RegisterHandler registers the user.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	decoder := json.NewDecoder(r.Body)

	var userRegisterForm forms.UserRegisterForm
	err := decoder.Decode(&userRegisterForm)
	if err != nil {
		error := make(map[string]string)
		error["message"] = "Invalid input format"
		response := Response{
			Success: false,
			Data:    nil,
			Error:   error,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	user := db.GetUserByEmail(userRegisterForm.Email)
	validationMessage := userRegisterForm.Validate(user)
	if validationMessage != "" {
		error := make(map[string]string)
		error["message"] = validationMessage
		response := Response{
			Success: false,
			Data:    nil,
			Error:   error,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	user := db.RegisterUser(userRegisterForm)

	objectID := db.GenerateObjectID()
	user.ID = objectID.Hex()
	db.CreateUser(user)
}

// LoginHandler validates the password & returns the JWT token.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	decoder := json.NewDecoder(r.Body)

	var userLoginForm forms.UserLoginForm
	err := decoder.Decode(&userLoginForm)
	if err != nil {
		panic(err)
	}
	user := db.GetUserByEmail(userLoginForm.Email)
	jwt := db.GetJWT(user, userLoginForm.Password)
	data := make(map[string]string)
	data["token"] = jwt

	response := Response{
		Success: true,
		Data:    data,
		Error:   nil,
	}
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

	user := db.GetUserByEmail(forgotPasswordForm.Email)
	toEmail := user.Email

	code := uuid.Must(uuid.NewV4())
	hexCode := hex.EncodeToString(code.Bytes())

	resetPassword := db.ResetPassword{UserID: user.ID, Code: hexCode}

	db.AddResetPassword(resetPassword)

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
	w.WriteHeader(http.StatusOK)

	decoder := json.NewDecoder(r.Body)
	var resetPasswordForm forms.ResetPasswordForm
	err := decoder.Decode(&resetPasswordForm)
	if err != nil {
		panic(err)
	}

	success := db.PasswordReset(resetPasswordForm.UID, resetPasswordForm.Code, resetPasswordForm.NewPassword)

	if success == true {
		user := db.GetUserByID(resetPasswordForm.UID)

		sub := "Forgot password"
		msg := fmt.Sprintf(
			"Hi," +
				"Your password has been reset successfully.",
		)

		utils.SendMail(sub, msg, user.Email)
	}
}
