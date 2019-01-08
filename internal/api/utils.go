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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

// This is to be used when the data has to be of the format
// Example:
// {
// 	  "data": {
// 		  "token": "abcde"
// 	  }
// }
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

// Use this when the value is complex data structure
// Example:
// {
// 	"data": {
// 		"user": {
// 			"firstName": "Shin",
// 			"lastName": "Lim"
// 		}
// 	}
// }
func writeSuccessStructResponse(w http.ResponseWriter, data map[string]interface{}, status int) {
	response := StructResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Use this when the data to be represented has no nested info.
// Example:
// {
// 	"data": {
// 		"firstName": "Shin",
// 		"lastName": "Lim"
// 	}
// }
func writeSuccessSimpleResponse(w http.ResponseWriter, data interface{}, status int) {
	response := SimpleResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
