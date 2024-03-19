package errresp

import (
	"github.com/groshi-project/groshi/pkg/httpresp"
	"net/http"
)

var InternalServerError = httpresp.New(
	http.StatusInternalServerError,
	NewErrorResponse("internal server error"),
)

var InvalidRequest = httpresp.New( // todo: rename InvalidRequest -> ?
	http.StatusBadRequest,
	NewErrorResponse("unable to decode request body"),
)

var InvalidRequestParams = httpresp.New(
	http.StatusBadRequest,
	NewErrorResponse("invalid request parameters"),
)

var UserNotFound = httpresp.New(
	http.StatusNotFound,
	NewErrorResponse("user not found"),
)

var CategoryNotFound = httpresp.New(
	http.StatusNotFound,
	NewErrorResponse("category not found"),
)

var CategoryForbidden = httpresp.New(
	http.StatusForbidden,
	NewErrorResponse("you have no access to this category"),
)

var InvalidCredentials = httpresp.New(
	http.StatusUnauthorized,
	NewErrorResponse("invalid credentials"),
)

var CurrencyNotFound = httpresp.New(
	http.StatusNotFound,
	NewErrorResponse("currency not found"),
)
