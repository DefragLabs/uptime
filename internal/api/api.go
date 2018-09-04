package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer Start the server.
func StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8080", router)
}
