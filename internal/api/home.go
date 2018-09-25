package api

import (
	"fmt"
	"net/http"
)

// HomeHandler - handler for root path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Home page")
}
