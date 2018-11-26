package api

import (
	"encoding/json"
	"net/http"
)

func writeErrorResponse(w http.ResponseWriter, errorMsg string) {
	errorVal := make(map[string]string)
	errorVal["message"] = errorMsg
	response := Response{
		Success: false,
		Data:    nil,
		Error:   errorVal,
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}
