package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler/httpresp"
	"github.com/groshi-project/groshi/internal/service/handler/response"
	"net/http"
	"time"
)

type transactionsCreateParams struct {
	Amount       int32  `json:"amount" example:"2500" validate:"required"`
	CurrencyCode string `json:"currency" example:"USD" validate:"required"`

	Timestamp time.Time `json:"timestamp" example:"todo-timestamp" validate:"required"`

	Description  string `json:"description" example:"Bought a donut for $2.5 only!"`
	CategoryUUID string `json:"category" example:"02983837-7ab0-492a-90b6-285491936067"`
}

type transactionsCreateResponse struct {
	UUID string `json:"uuid" example:"3be1ed0a-c307-49de-872e-38730200f301"`
}

// TransactionsCreate creates a new transactions and returns its UUID.
//
//	@Summary		Create a new transaction
//	@Description	Creates a new transaction and returns its UUID
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			user	body		transactionsCreateParams	true	"Transaction"
//	@Success		200		{object}	transactionsCreateResponse	"Successful operation"
//	@Failure		400		{object}	model.Error					"Invalid request body format or invalid request params"
//	@Failure		403		{object}	model.Error					"Access to the category is forbidden"
//	@Failure		404		{object}	model.Error					"User not found"
//	@Failure		500		{object}	model.Error					"Internal server error"
//	@Security		Bearer
//	@Router			/transactions [post]
func (h *Handler) TransactionsCreate(w http.ResponseWriter, r *http.Request) {
	// decode request params:
	params := &transactionsCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		httpresp.Render(w, response.InvalidRequestBodyFormat)
		return
	}

	// validate request params:
	if err := h.paramsValidate.Struct(params); err != nil {
		httpresp.Render(w, response.InvalidRequestParams)
		return
	}

	// fetch provided currency:
	currency := &database.Currency{}
	if err := h.database.SelectCurrencyByCode(r.Context(), params.CurrencyCode, currency); err != nil {
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
	if err := h.database.SelectCategoryByUUID(r.Context(), params.CategoryUUID, category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.CategoryNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// extract current user's username from context:
	username, ok := r.Context().Value("username").(string)
	if !ok {
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// fetch the current user from the database:
	user := &database.User{}
	if err := h.database.SelectUserByUsername(r.Context(), username, user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpresp.Render(w, response.UserNotFound)
			return
		}
		h.internalServerErrorLogger.Println(err)
		httpresp.Render(w, response.InternalServerError)
		return
	}

	// check if the given category belongs to the current user:
	if category.OwnerID != user.ID {
		httpresp.Render(w, response.CategoryForbidden)
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
	if err := h.database.CreateTransaction(r.Context(), transaction); err != nil {
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
	//uuid := chi.URLParam(r, "uuid")
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
