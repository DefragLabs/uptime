package api

import (
	"encoding/json"
	"net/http"

	"github.com/fatih/structs"
)

func writeErrorResponse(w http.ResponseWriter, errorMsg string) {
	errorVal := make(map[string]string)
	errorVal["message"] = errorMsg
	response := Response{
		Success: false,
		Data:    nil,
		Error:   errorVal,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

func writeSuccessResponse(w http.ResponseWriter, data map[string]string, status int) {
	response := Response{
		Success: true,
		Data:    data,
		Error:   nil,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func writeSuccessStructResponse(w http.ResponseWriter, data interface{}, status int) {
	responseData := structs.Map(data)
	response := StructResponse{
		Success: true,
		Data:    responseData,
		Error:   nil,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
