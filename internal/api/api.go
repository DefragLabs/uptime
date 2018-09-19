package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func monitoringDetailsRoutes(router *mux.Router) {
	router.HandleFunc("/monitoring-details", AddMonitoringDetailHandler).Methods("POST")
	router.HandleFunc("/monitoring-details", GetMonitoringDetailHandler).Methods("GET")
	router.HandleFunc("/monitoring-details", UpdateMonitoringDetailHandler).Methods("PUT")
}

// StartServer Start the server.
func StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/ping", PingHandler)

	monitoringDetailsRoutes(router)

	http.ListenAndServe(":8080", router)
}
