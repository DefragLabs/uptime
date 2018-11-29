package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
)

func clearUsersCollection() {
	datastore := db.New()
	datastore.Client.Database(datastore.DatabaseName).Collection(db.UsersCollection).Drop(context.Background())
}

func TestRegisterHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	defer clearUsersCollection()

	userRegisterForm := forms.UserRegisterForm{
		FirstName:   "Alice",
		LastName:    "Wonderland",
		Email:       "alice@sample.com",
		Password:    "test@123",
		CompanyName: "skynet",
	}
	byte, _ := json.Marshal(userRegisterForm)
	req, err := http.NewRequest("GET", "localhost:8080/api/auth/register", bytes.NewBuffer(byte))

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	RegisterHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status CREATED, got %v", res.StatusCode)
	}

	response := Response{}
	json.NewDecoder(res.Body).Decode(&response)

	if response.Success == false {
		t.Errorf("response success is false")
	} else if response.Data["token"] == "" {
		t.Errorf("token not found in the response")
	}
}
