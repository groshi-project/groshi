package errresp

import (
	"github.com/groshi-project/groshi/pkg/httpresp"
	"net/http"
)

var InternalServerError = httpresp.New(
	http.StatusInternalServerError,
	NewErrorData("internal server error"),
)

var InvalidRequest = httpresp.New( // todo: rename InvalidRequest -> ?
	http.StatusBadRequest,
	NewErrorData("unable to decode request body"),
)

var InvalidRequestParams = httpresp.New(
	http.StatusBadRequest,
	NewErrorData("invalid request parameters"),
)

var UserNotFound = httpresp.New(
	http.StatusNotFound,
	NewErrorData("user not found"),
)

var CategoryNotFound = httpresp.New(
	http.StatusNotFound,
	NewErrorData("category not found"),
)

var CategoryForbidden = httpresp.New(
	http.StatusForbidden,
	NewErrorData("you have no access to this category"),
)

var InvalidCredentials = httpresp.New(
	http.StatusUnauthorized,
	NewErrorData("invalid credentials"),
)

var CurrencyNotFound = httpresp.New(
	http.StatusNotFound,
	NewErrorData("currency not found"),
)
