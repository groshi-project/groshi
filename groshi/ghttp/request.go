package ghttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jieggii/groshi/groshi/handles/schema"
	"github.com/jieggii/groshi/groshi/loggers"
	"io"
	"net/http"
)

// Request is an abstraction which keeps request data
type Request struct {
	ResponseWriter http.ResponseWriter
	RawRequest     *http.Request
}

// sendJSONResponse sends JSON response.
func (req *Request) sendJSONResponse(data interface{}) {
	req.ResponseWriter.Header().Set("Content-Type", "application/json")
	req.ResponseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(req.ResponseWriter).Encode(data)
}

// Decode decodes request body.
func (req *Request) Decode(v interface{}) error {
	body, err := io.ReadAll(req.RawRequest.Body)
	if err != nil {
		panic(err) // todo
	}
	// todo:
	req.RawRequest.Body = io.NopCloser(bytes.NewBuffer(body))
	err = json.NewDecoder(req.RawRequest.Body).Decode(&v)
	req.RawRequest.Body = io.NopCloser(bytes.NewBuffer(body))

	return err
}

// DecodeSafe decodes request body and handles error.
// Returns true if there was no error.
func (req *Request) DecodeSafe(v interface{}) bool {
	if err := req.Decode(v); err != nil {
		fmt.Println("decode err", err)
		req.HandleError(err, schema.ClientSideError, schema.InvalidRequestBody)
		return false
	}
	return true
}

// HandleError handles any provided error. Does nothing if provided error is nil.
// Returns true if error was handled.
func (req *Request) HandleError(err error, errorCode schema.ErrorCode, errorMessage string) bool {
	if err == nil {
		return false
	}
	req.SendErrorResponse(errorCode, errorMessage, err)
	return true
}

// WrapCondition is wrapper for request body validations.
// Returns result value.
//func (req *Request) WrapCondition(result bool, errorMessage string) bool {
//	if !result {
//		req.SendErrorResponse(
//			schema.ClientSideError, errorMessage, nil,
//		)
//	}
//	return result
//}

// SendSuccessResponse sends success response.
func (req *Request) SendSuccessResponse(data interface{}) {
	successObject := schema.SuccessResponse{Success: true, Data: data}
	req.sendJSONResponse(successObject)
}

// SendErrorResponse sends error response. Sends "Internal server error"
// and logs error if it is server side error.
func (req *Request) SendErrorResponse(errorCode schema.ErrorCode, errorMessage string, err error) {
	if errorCode == schema.ServerSideError {
		loggers.Warn.Printf("Internal server error: `%v` (%v).\n", errorMessage, err)
		errorMessage = "Internal server error."
	}
	errorObject := schema.ErrorResponse{Success: false, ErrorCode: errorCode, ErrorMessage: errorMessage}
	req.sendJSONResponse(errorObject)
}

// NewRequest creates Request object.
func NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	return &Request{ResponseWriter: w, RawRequest: r}
}

// NewSafelyParsedRequest creates Request object and safely parses request body.
//func NewSafelyParsedRequest(w http.ResponseWriter, r *http.Request, v interface{}) (*Request, bool) {
//	req := NewRequest(w, r)
//	ok := req.DecodeSafe(v)
//	return req, ok
//}
