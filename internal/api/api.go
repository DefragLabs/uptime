package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func authRoutes(router *mux.Router) {
	s := router.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/login", LoginHandler).Methods("POST")
	s.HandleFunc("/logout", LogoutHandler).Methods("DELETE")
	s.HandleFunc("/register", RegisterHandler).Methods("POST")
	s.HandleFunc("/forgot-password", ForgotPasswordHandler).Methods("POST")
	s.HandleFunc("/reset-password", ResetPasswordHandler).Methods("POST")
}

func userRoutes(router *mux.Router) {
	router.HandleFunc("/users", GetUserDetailHandler).Methods("GET")
	router.HandleFunc("/users", UpdateUserDetailHandler).Methods("PUT")
}

func dashboardRoutes(router *mux.Router) {
	router.HandleFunc("/dashboard/stats", DashboardStatsHandler).Methods("GET")
}

func monitoringDetailsRoutes(router *mux.Router) {
	router.HandleFunc("/monitoring-urls", AddMonitoringURLHandler).Methods("POST")
	router.HandleFunc("/monitoring-urls", GetMonitoringURLsHandler).Methods("GET")

	// with search filter
	router.HandleFunc("/monitoring-urls", GetMonitoringURLsHandler).Queries(
		"search", "{search}",
	).Methods("GET")

	router.HandleFunc("/monitoring-urls/{monitoringURLID}", UpdateMonitoringURLHandler).Methods("PUT")
	router.HandleFunc("/monitoring-urls/{monitoringURLID}", GetMonitoringURLHandler).Methods("GET")
	router.HandleFunc("/monitoring-urls/{monitoringURLID}", DeleteMonitoringURLHandler).Methods("DELETE")
}

func monitoringStatsRoutes(router *mux.Router) {
	router.HandleFunc("/monitoring-urls/{monitoringURLID}/stats", GetMonitoringURLStatsHandler).Methods("GET")

	// with filters
	router.HandleFunc("/monitoring-urls/{monitoringURLID}/stats", GetMonitoringURLStatsHandler).Queries(
		"interval", "{interval}",
	).Methods("GET")
}

func integrationRoutes(router *mux.Router) {
	router.HandleFunc("/integrations", AddIntegrationHandler).Methods("POST")
	router.HandleFunc("/integrations", GetIntegrationsHandler).Methods("GET")
	router.HandleFunc("/integrations/{integrationID}", GetIntegrationHandler).Methods("GET")
	router.HandleFunc("/integrations/{integrationID}", DeleteIntegrationHandler).Methods("DELETE")
}

// StartServer Start the server.
func StartServer() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.HandleFunc("/", HomeHandler)

	monitoringDetailsRoutes(router)
	monitoringStatsRoutes(router)
	integrationRoutes(router)
	authRoutes(router)
	userRoutes(router)
	dashboardRoutes(router)

	http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}))(router),
	)
}
