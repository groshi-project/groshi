// Package response contains predefined HTTP error responses.
package response

import (
	"github.com/groshi-project/groshi/internal/service/handler/model"
	"github.com/groshi-project/groshi/pkg/httpresp"
	"net/http"
)

var InternalServerError = httpresp.New(
	http.StatusInternalServerError,
	model.NewError("internal server error"),
)

var InvalidRequest = httpresp.New( // todo: rename InvalidRequest -> ?
	http.StatusBadRequest,
	model.NewError("unable to decode request body"),
)

var InvalidRequestParams = httpresp.New(
	http.StatusBadRequest,
	model.NewError("invalid request parameters"),
)

var UserNotFound = httpresp.New(
	http.StatusNotFound,
	model.NewError("user not found"),
)

var CategoryNotFound = httpresp.New(
	http.StatusNotFound,
	model.NewError("category not found"),
)

var CategoryForbidden = httpresp.New(
	http.StatusForbidden,
	model.NewError("you have no access to this category"),
)

var InvalidCredentials = httpresp.New(
	http.StatusUnauthorized,
	model.NewError("invalid credentials"),
)

var CurrencyNotFound = httpresp.New(
	http.StatusNotFound,
	model.NewError("currency not found"),
)
