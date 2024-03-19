package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler/response"
	"github.com/groshi-project/groshi/pkg/httpresp"
	"net/http"
	"time"
)

type transactionsCreateParams struct {
	Amount       int32  `json:"amount"`
	CurrencyCode string `json:"currency"`

	Description  string `json:"description"`
	CategoryUUID string `json:"category"`

	Timestamp time.Time `json:"timestamp"`
}

type transactionsCreateResponse struct {
	UUID string `json:"uuid"`
}

func (h *Handler) TransactionsCreate(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &transactionsCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		httpresp.Render(w, response.InvalidRequest)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		httpresp.Render(w, response.InvalidRequestParams)
		return
	}

	// fetch provided currency:
	currency := &database.Currency{}
	if err := h.database.SelectCurrencyByCode(params.CurrencyCode, currency); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.CurrencyNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// fetch provided category:
	category := &database.Category{}
	if err := h.database.SelectCategoryByUUID(params.CategoryUUID, category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.CategoryNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// extract current user's username from claims:
	username, err := h.JWTAuthority.ExtractUsername(r.Context())
	if err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// fetch the current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.UserNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// convert provided timestamp to UTC timezone:
	utcTimestamp := params.Timestamp.UTC()

	// create a new transaction owned by the current user:
	transaction := &database.Transaction{
		Amount:     params.Amount,
		CurrencyID: currency.ID,

		Description: params.Description,
		CategoryID:  category.ID,

		OwnerID: user.ID,

		Timestamp: utcTimestamp,
	}
	if err := h.database.CreateTransaction(transaction); err != nil {
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// respond:
	resp := &transactionsCreateResponse{
		UUID: transaction.UUID.String(),
	}
	httpresp.NewOK(resp)
}

func (h *Handler) TransactionsGetOne(w http.ResponseWriter, r *http.Request) {

}

type transactionsGetParams struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Currency  string    `json:"currency"`
}

// todo: optionall by category, show transactions by categories
func (h *Handler) TransactionsGet(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) TransactionsUpdate(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) TransactionsDelete(w http.ResponseWriter, r *http.Request) {

}
