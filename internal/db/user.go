package db

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/dineshs91/uptime/internal/forms"
)

// User details struct
type User struct {
	ID           string `bson:"_id" json:"id,omitempty"`
	FirstName    string `bson:"firstName" json:"firstName"`
	LastName     string `bson:"lastName" json:"lastName"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"passwordHash" json:"password"`
}

// RegisterUser converts the password into hash and returns the user.
func RegisterUser(userRegisterForm forms.UserRegisterForm) User {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(userRegisterForm.Password), 14)

	user := User{
		FirstName:    userRegisterForm.FirstName,
		LastName:     userRegisterForm.LastName,
		Email:        userRegisterForm.Email,
		PasswordHash: string(passwordHash),
	}
	return user
}

// GetJWT validates user with password and returns JWT token.
func GetJWT(user User, password string) string {
	fmt.Println(password, user.PasswordHash, user)
	if checkPasswordHash(password, user.PasswordHash) == true {
		return "JWT success"
	}

	panic("Password check failed")
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
