package handles

import (
	"errors"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"github.com/jieggii/groshi/internal/http/handles/datatypes"
	"time"
)

type transactionCreateRequest struct {
	// required params:
	Amount *float64 `json:"amount"`

	// optional params:
	Description string                `json:"description"`
	Date        datatypes.ISO8601Date `json:"date"`
}

func (p *transactionCreateRequest) Before() error {
	if p.Amount == nil {
		return errors.New("missing required field `amount`")
	}

	// validate amount
	if *p.Amount <= 0 {
		return errors.New("`amount` must be more than 0")
	}

	// validate description (todo)

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
		Description: params.Description,
		Date:        params.Date.Time,

		OwnerId: currentUser.ID,
	}
	_, err := database.Db.NewInsert().Model(&transaction).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not create new transaction", err,
		)
		return
	}

	response := transactionCreateResponse{
		UUID: transaction.UUID,
	}
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

	Amount      float64 `json:"amount"`
	Description string  `json:"description"`

	Owner string    `json:"owner"`
	Date  time.Time `json:"date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func makeTransactionReadResponse(transaction *database.Transaction, transactionOwner *database.User) transactionReadResponse {
	return transactionReadResponse{
		UUID:        transaction.UUID,
		Amount:      transaction.Amount,
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

	NewAmount      *float64               `json:"new_amount"`
	NewDescription string                 `json:"new_description"`
	NewDate        *datatypes.ISO8601Date `json:"new_date"`
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
		transaction.Date = params.NewDate.Time
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
	Since datatypes.ISO8601Date `json:"since"`

	// optional options
	Until datatypes.ISO8601Date `json:"until"`
}

func (p *transactionListRequest) Before() error {
	// validate since
	if !p.Since.IsValid {
		return errors.New("invalid value of field `since`")
	}

	// validate until
	if !p.Until.IsValid {
		return errors.New("invalid value of field `until`")
	}

	if p.Since.Time.IsZero() {
		return errors.New("missing required field `since`")
	}

	if p.Until.Time.IsZero() {
		p.Until = datatypes.ISO8601Date{Time: time.Now()}
	}

	return nil
}

type transactionListResponse struct {
	Count        int                       `json:"count"`
	Transactions []transactionReadResponse `json:"transactions"` // todo: do something with problem: null is returned instead of [] when Count == 0.
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

type transactionSummaryRequest struct {
	// required options
	Since datatypes.ISO8601Date `json:"since"`

	// optional options
	Until datatypes.ISO8601Date `json:"until"`
}

func (p *transactionSummaryRequest) Before() error {
	// validate since
	if !p.Since.IsValid {
		return errors.New("invalid value of field `since`")
	}

	// validate until
	if !p.Until.IsValid {
		return errors.New("invalid value of field `until`")
	}

	if p.Since.Time.IsZero() {
		return errors.New("missing required field `since`")
	}
	if p.Until.Time.IsZero() {
		p.Until = datatypes.ISO8601Date{Time: time.Now()}
	}
	return nil
}

type transactionSummaryResponse struct {
	Count   int     `json:"count"`
	Income  float64 `json:"income"`
	Outcome float64 `json:"outcome"`
	Total   float64 `json:"total"`
}

// TransactionSummary returns summary for transactions owned by current user for given period.
// (count, income, outcome, total)
func TransactionSummary(request *ghttp.Request, currentUser *database.User) {
	params := transactionSummaryRequest{}
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
		).Scan(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not fetch transactions", err,
		)
		return
	}

	income := 0.0
	outcome := 0.0
	total := 0.0

	for _, transaction := range transactions {

		if transaction.Amount >= 0 {
			income += transaction.Amount
		} else {
			outcome += transaction.Amount
		}
		total += transaction.Amount
	}

	response := transactionSummaryResponse{
		Count:   len(transactions),
		Income:  income,
		Outcome: outcome,
		Total:   total,
	}
	request.SendSuccessfulResponse(&response)
}
