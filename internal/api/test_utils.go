package api

import (
	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
)

func createTestUser() (db.User, string) {
	userRegisterForm := forms.UserRegisterForm{
		FirstName:   "Alice",
		LastName:    "Wonderland",
		Email:       "alice@sample.com",
		Password:    "test@123",
		CompanyName: "skynet",
	}
	newUser := db.RegisterUser(userRegisterForm)
	objectID := db.GenerateObjectID()
	newUser.ID = objectID.Hex()

	datastore := db.New()
	datastore.CreateUser(newUser)
	jwt, _ := db.GetJWT(newUser, userRegisterForm.Password)
	return newUser, jwt
}
