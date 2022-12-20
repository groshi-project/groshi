package util

import (
	"encoding/json"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/loggers"
	"net/http"
)

func returnJSON(writer http.ResponseWriter, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(object)
}

func DecodeBodyJSON(writer http.ResponseWriter, request *http.Request, object interface{}) bool {
	err := json.NewDecoder(request.Body).Decode(object)
	if err != nil {
		ReturnErrorResponse(writer, schema.ClientSideError, "Invalid request body.", nil)
		return false
	}
	return true
}

func ReturnSuccessResponse(writer http.ResponseWriter, data interface{}) {
	successObject := schema.SuccessResponse{Success: true, Data: data}
	returnJSON(writer, successObject)
}

func ReturnErrorResponse(
	writer http.ResponseWriter,
	errorCode schema.ErrorCode,
	errorMessage string,
	err error,
) {
	if errorCode == schema.ServerSideError {
		loggers.Warn.Printf("Internal server error: `%v` (%v).\n", errorMessage, err)
		errorMessage = "Internal server error."
	}
	errorObject := schema.ErrorResponse{Success: false, ErrorCode: errorCode, ErrorMessage: errorMessage}
	returnJSON(writer, errorObject)
}
