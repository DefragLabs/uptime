package api

import (
	"encoding/json"
	"net/http"

	"github.com/dineshs91/uptime/internal/db"
	"github.com/dineshs91/uptime/internal/forms"
)

// RegisterHandler registers the user.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	decoder := json.NewDecoder(r.Body)

	var userRegisterForm forms.UserRegisterForm
	err := decoder.Decode(&userRegisterForm)
	if err != nil {
		panic(err)
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
	user := db.GetUser(userLoginForm.Email)
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
