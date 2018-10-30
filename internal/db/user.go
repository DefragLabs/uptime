package db

import (
	"os"
	"time"

	"github.com/defraglabs/uptime/internal/forms"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User details struct
type User struct {
	ID           string `bson:"_id" json:"id,omitempty"`
	FirstName    string `bson:"firstName" json:"firstName"`
	LastName     string `bson:"lastName" json:"lastName"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"passwordHash" json:"password"`
}

// ResetPassword struct
type ResetPassword struct {
	ID     string `bson:"_id" json:"id,omitempty"`
	UserID string `bson:"userID"`
	Code   string `bson:"code"`
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
	if checkPasswordHash(password, user.PasswordHash) == true {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 168).Unix(),
			"iat":   time.Now().Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		return tokenString
	}

	panic("Password check failed")
}

// PasswordReset validates & resets the password.
func PasswordReset(uid, code, newPassword string) bool {
	datastore := New()
	user := datastore.GetUserByID(uid)
	resetPassword := datastore.GetResetPassword(uid, code)

	if resetPassword.ID != "" {
		// updatePassword.
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
		user.PasswordHash = string(passwordHash)
		// UpdateUser(user)
		return true
	}

	return false
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
