package handlers

import (
	"encoding/json"
	"net/http"
)

func decodeBodyJSON(writer http.ResponseWriter, request *http.Request, object interface{}) bool {
	err := json.NewDecoder(request.Body).Decode(object)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

func returnJSON(writer http.ResponseWriter, httpStatusCode int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatusCode)
	json.NewEncoder(writer).Encode(object)
}

func returnError(writer http.ResponseWriter, httpStatusCode int, errorMessage string) {
	errorObject := ErrorResponse{ErrorMessage: errorMessage}
	returnJSON(writer, httpStatusCode, errorObject)
}
