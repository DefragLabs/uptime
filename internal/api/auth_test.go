package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

	req, err := http.NewRequest("GET", "localhost:8080/register", strings.NewReader(fmt.Sprintf("%#v", userRegisterForm)))

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	RegisterHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	print(res.StatusCode)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status CREATED, got %v", res.StatusCode)
	}
}
