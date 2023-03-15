// Package ghttp (stands for "groshi HTTP") is basically a tiny HTTP framework
// which simplifies reading HTTP request and writing HTTP responses.
package ghttp

import (
	"bytes"
	"encoding/json"
	"github.com/jieggii/groshi/internal/database/currency"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"github.com/jieggii/groshi/internal/loggers"
	"io"
	"net/http"
)

// Request is struct which keeps request data.
type Request struct {
	ResponseWriter http.ResponseWriter
	RawRequest     *http.Request
}

// sendJSONResponse sends JSON response.
func (req *Request) sendJSONResponse(data interface{}) {
	req.ResponseWriter.Header().Set("Content-Type", "application/json")
	req.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*") // enable CORS
	req.ResponseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(req.ResponseWriter).Encode(data)
}

// Decode decodes request body and handles all possible errors.
// Returns true if there was no error.
func (req *Request) Decode(params RequestParams) bool {
	body, err := io.ReadAll(req.RawRequest.Body)
	if err != nil {
		req.SendServerSideErrorResponse("could not read request body", err)
		return false
	}

	req.RawRequest.Body = io.NopCloser(bytes.NewBuffer(body))
	err = json.NewDecoder(req.RawRequest.Body).Decode(params)
	req.RawRequest.Body = io.NopCloser(bytes.NewBuffer(body))

	if err != nil {
		switch err.(type) {
		case *currency.UnknownCurrencyError: // could not unmarshal currency
			req.SendClientSideErrorResponse(
				schema.InvalidRequestErrorTag,
				"unknown currency",
			)
		default:
			req.SendClientSideErrorResponse(
				schema.InvalidRequestErrorTag,
				"could not parse request (probably incorrect format or type of some fields)",
			)
		}

		return false
	}
	return true
}

// SendSuccessfulResponse sends successful response.
func (req *Request) SendSuccessfulResponse(data Response) {
	successObject := schema.SuccessResponse{Success: true, Data: data}
	req.sendJSONResponse(successObject)
}

// SendClientSideErrorResponse sends client-side error response containing error message.
func (req *Request) SendClientSideErrorResponse(errorTag schema.ErrorTag, errorDetails string) {
	req.sendJSONResponse(schema.ErrorResponse{
		Success:      false,
		ErrorOrigin:  schema.ClientErrorOrigin,
		ErrorTag:     errorTag,
		ErrorDetails: errorDetails,
	})
}

// SendServerSideErrorResponse sends server-side error response without details and logs the error.
func (req *Request) SendServerSideErrorResponse(errorComment string, err error) {
	req.sendJSONResponse(schema.ErrorResponse{
		Success:      false,
		ErrorOrigin:  schema.ServerErrorOrigin,
		ErrorTag:     schema.InternalServerErrorErrorTag,
		ErrorDetails: "Internal server error.",
	})
	loggers.Error.Printf("returned server-side error response 'cause %v (%v)", errorComment, err)
}

// NewRequest creates new Request instance.
func NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	return &Request{ResponseWriter: w, RawRequest: r}
}
