package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func writeJsonErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	response := errorResponse{
		Message: err.Error(),
	}
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
