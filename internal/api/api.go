package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer Start the server.
func StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/ping", PingHandler)
	http.ListenAndServe(":8080", router)
}
