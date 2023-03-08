package ghttp

import (
	"bytes"
	"encoding/json"
	"github.com/jieggii/groshi/internal/loggers"
	schema2 "github.com/jieggii/groshi/internal/schema"
	"io"
	"net/http"
)

// Request is an abstraction which keeps request data
type Request struct {
	ResponseWriter http.ResponseWriter
	RawRequest     *http.Request
}

// NewRequest creates new Request instance.
func NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	return &Request{ResponseWriter: w, RawRequest: r}
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
// todo?: merge Decode and DecodeSafe
func (req *Request) DecodeSafe(v interface{}) bool {
	if err := req.Decode(v); err != nil {
		req.SendClientSideErrorResponse(schema2.InvalidRequestBody)
		return false
	}
	return true
}

// SendSuccessResponse sends success response.
func (req *Request) SendSuccessResponse(data interface{}) {
	successObject := schema2.SuccessResponse{Success: true, Data: data}
	req.sendJSONResponse(successObject)
}

// SendClientSideErrorResponse sends client-side error response containing error message.
func (req *Request) SendClientSideErrorResponse(errorMessage string) {
	req.sendJSONResponse(schema2.ErrorResponse{
		Success:      false,
		ErrorOrigin:  schema2.ErrorOriginClient,
		ErrorMessage: errorMessage,
	})
}

// SendServerSideErrorResponse sends server-side error response containing error message
// "Internal Server Error", logs error comment and error object.
func (req *Request) SendServerSideErrorResponse(errorComment string, err error) {
	req.sendJSONResponse(schema2.ErrorResponse{
		Success:      false,
		ErrorOrigin:  schema2.ErrorOriginServer,
		ErrorMessage: schema2.InternalServerError,
	})
	loggers.Error.Printf("%v (%v)", errorComment, err)
}
