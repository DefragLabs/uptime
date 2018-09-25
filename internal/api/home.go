package api

import (
	"fmt"
	"net/http"

	"github.com/dineshs91/uptime/internal/tasks"
)

// HomeHandler - handler for root path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Home page")
}

// PingHandler - handler for ping path
func PingHandler(w http.ResponseWriter, r *http.Request) {
	go uptime.StartScheduler()
	fmt.Fprintf(w, "Pinging \n")
}
