package util

import (
	"encoding/json"
	"net/http"
)

func DecodeBodyJSON(writer http.ResponseWriter, request *http.Request, object interface{}) bool {
	err := json.NewDecoder(request.Body).Decode(object)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

func ReturnJSON(writer http.ResponseWriter, httpStatusCode int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatusCode)
	json.NewEncoder(writer).Encode(object)
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func ReturnError(writer http.ResponseWriter, httpStatusCode int, errorMessage string) {
	errorObject := ErrorResponse{ErrorMessage: errorMessage}
	ReturnJSON(writer, httpStatusCode, errorObject)
}
