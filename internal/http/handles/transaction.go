package handles

import (
	"errors"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/database/currency"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"time"
)

type transactionCreateRequest struct {
	// Required params:
	Amount   *float64          `json:"amount"`
	Currency currency.Currency `json:"currency"`

	// Optional params:
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (p *transactionCreateRequest) Before() error {
	if p.Amount == nil || p.Currency == "" {
		return errors.New("missing required fields `amount` and `currency`")
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

	if err := params.Before(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	transaction := database.Transaction{
		Amount:      *params.Amount,
		Currency:    params.Currency,
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

func (p *transactionReadRequest) Before() error {
	if p.UUID == "" {
		return errors.New("missing required field `uuid`")
	}
	return nil
}

type transactionReadResponse struct {
	UUID string `json:"uuid"`

	Amount      float64           `json:"amount"`
	Currency    currency.Currency `json:"currency"`
	Description string            `json:"description"`

	Owner string    `json:"owner"`
	Date  time.Time `json:"date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func makeTransactionReadResponse(transaction *database.Transaction, transactionOwner *database.User) transactionReadResponse {
	return transactionReadResponse{
		UUID:        transaction.UUID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,

		Owner: transactionOwner.Username,
		Date:  transaction.Date,

		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}
}

// TransactionRead returns information about transaction owned by current user.
func TransactionRead(request *ghttp.Request, currentUser *database.User) {
	params := transactionReadRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Before(); err != nil {
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

	response := makeTransactionReadResponse(transaction, currentUser)
	request.SendSuccessfulResponse(&response)
}

type transactionUpdateRequest struct {
	UUID string `json:"uuid"`

	NewAmount      *float64   `json:"new_amount"`
	NewDescription string     `json:"new_description"`
	NewDate        *time.Time `json:"new_date"`
}

func (p *transactionUpdateRequest) Before() error {
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

	if err := params.Before(); err != nil {
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

	_, err = database.Db.NewUpdate().Model(transaction).WherePK().Exec(database.Ctx)
	if err != nil {
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

func (p *transactionDeleteRequest) Before() error {
	if p.UUID == "" {
		return errors.New("missing required field `uuid`")
	}
	return nil
}

type transactionDeleteResponse struct {
	UUID string `json:"uuid"`
}

// TransactionDelete deletes transaction.
func TransactionDelete(request *ghttp.Request, currentUser *database.User) {
	params := transactionDeleteRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}

	if err := params.Before(); err != nil {
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

	response := transactionDeleteResponse{
		UUID: transaction.UUID,
	}
	request.SendSuccessfulResponse(&response)
}

type transactionListRequest struct {
	// required options
	Since time.Time `json:"since"`

	// optional options
	Until time.Time `json:"until"`
}

func (p *transactionListRequest) Before() error {
	if p.Since.IsZero() {
		return errors.New("missing required field `since`")
	}
	if p.Until.IsZero() {
		p.Until = time.Now()
	}
	return nil
}

type transactionListResponse struct {
	Count        int                       `json:"count"`
	Transactions []transactionReadResponse `json:"transactions"`
}

// TransactionList gets list of transactions owned by current user.
func TransactionList(request *ghttp.Request, currentUser *database.User) {
	params := transactionListRequest{}
	if ok := request.Decode(&params); !ok {
		return
	}
	if err := params.Before(); err != nil {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, err.Error(),
		)
		return
	}

	var transactions []database.Transaction
	err := database.Db.NewSelect().Model(&transactions).
		Where(
			"(owner_id = ?) AND (created_at BETWEEN ? and ?)",
			currentUser.ID, params.Since, params.Until,
		).Scan(database.Ctx, &transactions)

	if err != nil {
		request.SendServerSideErrorResponse(
			"could not fetch transactions", err,
		)
		return
	}

	var responseTransactions []transactionReadResponse

	for _, transaction := range transactions {
		responseTransactions = append(
			responseTransactions, makeTransactionReadResponse(&transaction, currentUser),
		)
	}

	response := transactionListResponse{
		Count:        len(transactions),
		Transactions: responseTransactions,
	}

	request.SendSuccessfulResponse(&response)
}
