package util

import (
	"encoding/json"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/loggers"
	"net/http"
)

func returnJSON(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(object)
}

func DecodeBodyJSON(w http.ResponseWriter, r *http.Request, object interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(object)
	if err != nil {
		ReturnErrorResponse(w, schema.ClientSideError, "Invalid request body.", nil)
		return false
	}
	return true
}

func ValidateBody(writer http.ResponseWriter, expression bool) bool {
	if !expression {
		ReturnErrorResponse(writer, schema.ClientSideError, "Invalid request body.", nil)
		return false
	}
	return true
}

func ReturnSuccessResponse(w http.ResponseWriter, data interface{}) {
	successObject := schema.SuccessResponse{Success: true, Data: data}
	returnJSON(w, successObject)
}

func ReturnErrorResponse(w http.ResponseWriter, errorCode schema.ErrorCode, errorMessage string, err error) {
	if errorCode == schema.ServerSideError {
		loggers.Warn.Printf("Internal server error: `%v` (%v).\n", errorMessage, err)
		errorMessage = "Internal server error."
	}
	errorObject := schema.ErrorResponse{Success: false, ErrorCode: errorCode, ErrorMessage: errorMessage}
	returnJSON(w, errorObject)
}
