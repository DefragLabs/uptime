package db

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/defraglabs/uptime/internal/cache"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User details struct
type User struct {
	ID           string `bson:"_id" json:"id,omitempty"`
	FirstName    string `bson:"firstName" json:"firstName"`
	LastName     string `bson:"lastName" json:"lastName"`
	PhoneNumber  string `bson:"phoneNumber" json:"phoneNumber"`
	CompanyName  string `bson:"CompanyName" json:"CompanyName"`
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
		PhoneNumber:  userRegisterForm.PhoneNumber,
		CompanyName:  userRegisterForm.CompanyName,
	}
	return user
}

// GetJWT validates user with password and returns JWT token.
func GetJWT(user User, password string) (string, error) {
	if checkPasswordHash(password, user.PasswordHash) == true {
		code := uuid.Must(uuid.NewV4())
		hexCode := hex.EncodeToString(code.Bytes())

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": user.ID,
			"email":  user.Email,
			"exp":    time.Now().Add(time.Hour * 168).Unix(),
			"iat":    time.Now().Unix(),
			"jti":    hexCode,
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		return tokenString, nil
	}

	return "", errors.New("invalid credentials")
}

func splitAuthToken(authToken string) ([]string, error) {
	splitAuthToken := strings.Split(authToken, " ")

	if len(splitAuthToken) != 2 {
		var emptyArray []string
		return emptyArray, errors.New("invalid auth header")
	}

	return splitAuthToken, nil
}

// ValidateJWT validates the provided JWT and returns the corresponding user.
func ValidateJWT(authToken string) (User, error) {
	authTokenArray, splitErr := splitAuthToken(authToken)

	if splitErr != nil {
		return User{}, errors.New(splitErr.Error())
	}

	authType, tokenString := authTokenArray[0], authTokenArray[1]

	if authType != "JWT" {
		return User{}, errors.New("invalid auth type")
	}

	claims, jwtDecodeErr := verifyJWT(tokenString)

	if jwtDecodeErr != nil {
		return User{}, errors.New(jwtDecodeErr.Error())
	}

	jti := claims["jti"].(string)

	if isTokenRevoked(jti) {
		return User{}, errors.New("token has been revoked")
	}
	userID := claims["userID"].(string)

	if userID != "" {
		datastore := New()
		user := datastore.GetUserByID(userID)

		if user.ID == "" {
			return User{}, errors.New("user not found")
		}

		return user, nil
	}

	return User{}, errors.New("authorization failed")
}

func isTokenRevoked(jti string) bool {
	cacheClient, cacheErr := cache.GetClient()

	if cacheErr != nil {
		return true
	}

	val, err := cacheClient.Get(jti).Result()

	if err != nil || val != "revoked" {
		return false
	}
	return true
}

// RevokeToken adds the token to redis blacklist.
// The token expires from the cache at `exp` + int32 of the token.
func RevokeToken(authToken string) bool {
	authTokenArray, _ := splitAuthToken(authToken)

	tokenString := authTokenArray[1]
	claims, tokenDecodeErr := verifyJWT(tokenString)

	if tokenDecodeErr != nil {
		return false
	}

	exp := int64(claims["exp"].(float64))
	exp += 1000

	jti := claims["jti"].(string)

	cacheClient, cacheErr := cache.GetClient()

	if cacheErr != nil {
		return false
	}

	duration := exp - time.Now().Unix()
	cacheClient.Set(jti, "revoked", time.Duration(duration)*time.Second)

	return true
}

// verifyJWT verifies the token & returns the payload.
func verifyJWT(tokenString string) (jwt.MapClaims, error) {
	hmacSampleSecret := []byte(os.Getenv("JWT_SECRET"))
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Info(nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"]))
		}

		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	log.Warn("Unable to decode jwt token")
	return jwt.MapClaims{}, errors.New("Unable to decode token")
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
