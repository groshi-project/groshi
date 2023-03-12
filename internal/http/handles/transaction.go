package handles

import (
	"errors"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"time"
)

type transactionCreateRequest struct {
	// Required params:
	Amount   *float64 `json:"amount"`
	Currency string   `json:"currency"`

	// Optional params:
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (p *transactionCreateRequest) Validate() error {
	if p.Amount == nil || p.Currency == "" {
		return errors.New("these fields are required: `amount`, `currency`, `description`")
	}

	return nil
}

type transactionCreateResponse struct {
	UUID string `json:"uuid"`
}

// TransactionCreate creates new transaction.
func TransactionCreate(request *ghttp.Request, currentUser *database.User) {
	params := transactionCreateRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Validate(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	ok, currency := database.Currencies.GetCurrency(params.Currency)
	if !ok {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.UnknownCurrencyErrorDetail,
		)
		return
	}

	transaction := database.Transaction{
		Amount:      *params.Amount,
		Currency:    currency,
		Description: params.Description,
		Date:        params.Date,

		OwnerId: currentUser.ID,
	}
	_, err := database.Db.NewInsert().Model(&transaction).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not create new transaction", err,
		)
		return
	}

	response := transactionCreateResponse{UUID: transaction.UUID}
	request.SendSuccessfulResponse(&response)
}

type transactionReadRequest struct {
	UUID string `json:"uuid"`
}

func (p *transactionReadRequest) Validate() error {
	if p.UUID == "" {
		return errors.New("missing required field `uuid`")
	}
	return nil
}

type transactionReadResponse struct {
	UUID string `json:"uuid"`

	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`

	Owner string    `json:"owner"`
	Date  time.Time `json:"date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// TransactionRead returns information about transaction.
func TransactionRead(request *ghttp.Request, currentUser *database.User) {
	params := transactionReadRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Validate(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	transaction, err := database.FetchTransactionByUUID(params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.TransactionNotFoundErrorDetail,
		)
		return
	}

	transactionOwner := database.User{}
	err = database.Db.NewSelect().Model(&transactionOwner).Where("id = ?", transaction.OwnerId).Scan(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not fetch transaction owner", err,
		)
		return
	}

	if transactionOwner.ID != currentUser.ID {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.ThisTransactionDoesNotBelongToYouErrorDetail,
		)
		return
	}

	response := transactionReadResponse{
		UUID:        transaction.UUID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,

		Owner: transactionOwner.Username,
		Date:  transaction.Date,

		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}
	request.SendSuccessfulResponse(&response)
}

type transactionUpdateRequest struct {
	UUID string `json:"uuid"`

	NewAmount      *float64   `json:"new_amount"`
	NewDescription string     `json:"new_description"`
	NewDate        *time.Time `json:"new_date"`
}

func (p *transactionUpdateRequest) Validate() error {
	if p.UUID == "" {
		return errors.New("missing required field `uuid`")
	}
	if p.NewAmount == nil && p.NewDescription == "" && p.NewDate == nil {
		return errors.New(
			"at least one of these fields is required: `new_amount`, `new_description`, `new_date`",
		)
	}
	return nil
}

type transactionUpdateResponse struct {
	UUID string `json:"uuid"`
}

// TransactionUpdate updates transaction.
func TransactionUpdate(request *ghttp.Request, currentUser *database.User) {
	params := transactionUpdateRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Validate(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	transaction, err := database.FetchTransactionByUUID(params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.TransactionNotFoundErrorDetail,
		)
		return
	}

	if transaction.OwnerId != currentUser.ID {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.ThisTransactionDoesNotBelongToYouErrorDetail,
		)
		return
	}

	if params.NewAmount != nil {
		transaction.Amount = *params.NewAmount
	}

	if params.NewDescription != "" {
		transaction.Description = params.NewDescription
	}

	if params.NewDate != nil {
		transaction.Date = *params.NewDate
	}

	if _, err := database.Db.NewUpdate().Model(transaction).WherePK().Exec(database.Ctx); err != nil {
		request.SendServerSideErrorResponse("could not update transaction", err)
		return
	}

	response := transactionUpdateResponse{
		UUID: transaction.UUID,
	}
	request.SendSuccessfulResponse(&response)
}

type transactionDeleteRequest struct {
	UUID string `json:"uuid"`
}

func (p *transactionDeleteRequest) Validate() error {
	if p.UUID == "" {
		return errors.New("missing required field `uuid`")
	}
	return nil
}

//type transactionDeleteResponse struct{}

// TransactionDelete deletes transaction.
func TransactionDelete(request *ghttp.Request, currentUser *database.User) {
	params := transactionDeleteRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Validate(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	transaction, err := database.FetchTransactionByUUID(params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.TransactionNotFoundErrorDetail,
		)
		return
	}

	if transaction.OwnerId != currentUser.ID {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.ThisTransactionDoesNotBelongToYouErrorDetail,
		)
		return
	}

	_, err = database.Db.NewDelete().Model(transaction).Where("uuid = ?", transaction.UUID).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("could not delete transaction", err)
		return
	}

	//response := transactionDeleteResponse{}
	request.SendSuccessfulResponse(&ghttp.EmptyResponse{})
}
