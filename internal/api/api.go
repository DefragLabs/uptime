package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func authRoutes(router *mux.Router) {
	s := router.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/login", LoginHandler).Methods("POST")
	s.HandleFunc("/register", RegisterHandler).Methods("POST")
	s.HandleFunc("/forgot-password", ForgotPasswordHandler).Methods("POST")
	s.HandleFunc("/reset-password", ResetPasswordHandler).Methods("POST")
}

func monitoringDetailsRoutes(router *mux.Router) {
	router.HandleFunc("/monitoring-urls", AddMonitoringURLHandler).Methods("POST")
	router.HandleFunc("/monitoring-urls", GetMonitoringURLHandler).Methods("GET")
	router.HandleFunc("/monitoring-urls", UpdateMonitoringURLHandler).Methods("PUT")
}

// StartServer Start the server.
func StartServer() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.HandleFunc("/", HomeHandler)

	monitoringDetailsRoutes(router)
	authRoutes(router)

	http.ListenAndServe(":8080", router)
}
