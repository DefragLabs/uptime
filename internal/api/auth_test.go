package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/defraglabs/uptime/internal/forms"
)

func TestRegisterHandler(t *testing.T) {
	userRegisterForm := forms.UserRegisterForm{
		FirstName: "Alice",
		LastName:  "Wonderland",
		Email:     "alice@sample.com",
		Password:  "test@123",
	}
	byte, _ := json.Marshal(userRegisterForm)
	req, err := http.NewRequest("GET", "localhost:8080/register", bytes.NewBuffer(byte))

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	RegisterHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	t.Errorf(string(b))
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status CREATED, got %v", res.StatusCode)
	}
}
