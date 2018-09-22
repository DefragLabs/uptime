package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func monitoringDetailsRoutes(router *mux.Router) {
	router.HandleFunc("/monitoring-urls", AddMonitoringURLHandler).Methods("POST")
	router.HandleFunc("/monitoring-urls", GetMonitoringURLHandler).Methods("GET")
	router.HandleFunc("/monitoring-urls", UpdateMonitoringURLHandler).Methods("PUT")
}

// StartServer Start the server.
func StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/ping", PingHandler)

	monitoringDetailsRoutes(router)

	http.ListenAndServe(":8080", router)
}
