package handles

import (
	"errors"
	"fmt"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"github.com/jieggii/groshi/internal/http/handles/datatypes"
	"github.com/jieggii/groshi/internal/http/handles/validators"
	"time"
)

type transactionCreateRequest struct {
	Amount   *float64            `json:"amount"`
	Currency *datatypes.Currency `json:"currency"`

	BaseAmount  *float64               `json:"base_amount"`
	Description *string                `json:"description"`
	Date        *datatypes.ISO8601Date `json:"date"`
}

func (p *transactionCreateRequest) Before() error {
	if p.Amount == nil || p.Currency == nil {
		return errors.New(
			schema.MissingRequiredFieldsErrorDetail("amount", "currency"),
		)
	}

	// validate currency
	if err := validators.ValidateCurrency(*p.Currency); err != nil {
		return err
	}

	if p.Description != nil {
		if err := validators.ValidateTransactionDescription(*p.Description); err != nil {
			return err
		}
	}

	if p.Date == nil {
		p.Date = &datatypes.ISO8601Date{Time: time.Now()}
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

	fmt.Println(currentUser.BaseCurrency)

	if currentUser.BaseCurrency != (*params.Currency).String() {
		if params.BaseAmount == nil {
			request.SendClientSideErrorResponse(
				schema.InvalidRequestErrorTag,
				"missing required field `base_amount` (this field is required because you are creating transaction with currency other than your base currency)",
			)
			return
		}
	}

	if params.BaseAmount == nil {
		params.BaseAmount = params.Amount
	}

	transaction := database.Transaction{
		Date:        (*params.Date).Time,
		OwnerId:     currentUser.ID,
		BaseAmount:  *params.BaseAmount,
		Amount:      *params.Amount,
		Currency:    *params.Currency,
		Description: params.Description,
	}

	_, err := database.Db.NewInsert().Model(&transaction).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"could not insert new transaction", err,
		)
		return
	}
	response := transactionCreateResponse{
		UUID: transaction.UUID,
	}
	request.SendSuccessfulResponse(&response)
}

type transactionReadRequest struct {
	UUID *string `json:"uuid"`
}

func (p *transactionReadRequest) Before() error {
	if p.UUID == nil {
		return errors.New(
			schema.MissingRequiredFieldErrorDetail("uuid"),
		)
	}
	return nil
}

type transactionReadResponse struct {
	UUID string `json:"uuid"`

	BaseAmount float64            `json:"base_amount"`
	Amount     float64            `json:"amount"`
	Currency   datatypes.Currency `json:"currency"`

	Description *string `json:"description"`

	Owner string    `json:"owner"`
	Date  time.Time `json:"date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func makeTransactionReadResponse(transaction *database.Transaction, transactionOwner *database.User) transactionReadResponse {
	return transactionReadResponse{
		UUID: transaction.UUID,

		BaseAmount: transaction.BaseAmount,
		Amount:     transaction.Amount,
		Currency:   transaction.Currency,

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

	transaction, err := database.GetTransaction(*params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.TransactionNotFoundErrorDetail,
		)
		return
	}

	if transaction.OwnerId != currentUser.ID {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.TransactionDoesNotBelongToYouErrorDetail,
		)
		return
	}

	response := makeTransactionReadResponse(transaction, currentUser)
	request.SendSuccessfulResponse(&response)
}

type transactionUpdateRequest struct {
	UUID *string `json:"uuid"`

	NewBaseAmount  *float64               `json:"new_base_amount"`
	NewAmount      *float64               `json:"new_amount"`
	NewDescription *string                `json:"new_description"`
	NewDate        *datatypes.ISO8601Date `json:"new_date"`
}

func (p *transactionUpdateRequest) Before() error {
	if p.UUID == nil {
		return errors.New(
			schema.MissingRequiredFieldErrorDetail("uuid"),
		)
	}

	if p.NewAmount == nil && p.NewDescription == nil && p.NewDate == nil {
		return errors.New(
			schema.AtLeastOneOfFieldsIsRequiredErrorDetail(
				"new_amount", "new_description", "new_date",
			),
		)
	}

	// validate date
	if p.NewDate != nil {
		if !p.NewDate.IsValid {
			return errors.New("invalid value of field `new_date`")
		}
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

	transaction, err := database.GetTransaction(*params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.TransactionNotFoundErrorDetail,
		)
		return
	}

	if transaction.OwnerId != currentUser.ID {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.TransactionDoesNotBelongToYouErrorDetail,
		)
		return
	}

	if params.NewAmount != nil {
		transaction.Amount = *params.NewAmount
	}

	if params.NewDescription != nil {
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
	UUID *string `json:"uuid"`
}

func (p *transactionDeleteRequest) Before() error {
	if p.UUID == nil {
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

	transaction, err := database.GetTransaction(*params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(
			schema.ObjectNotFoundErrorTag, schema.TransactionNotFoundErrorDetail,
		)
		return
	}

	if transaction.OwnerId != currentUser.ID {
		request.SendClientSideErrorResponse(
			schema.AccessDeniedErrorTag, schema.TransactionDoesNotBelongToYouErrorDetail,
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
	Since *datatypes.ISO8601Date `json:"since"`

	// optional options
	Until *datatypes.ISO8601Date `json:"until"`
}

func (p *transactionListRequest) Before() error {
	if p.Since == nil {
		return errors.New("missing required field `since`")
	}

	if p.Until == nil {
		p.Until = &datatypes.ISO8601Date{Time: time.Now()}
	}

	// validate since
	if !p.Since.IsValid {
		return errors.New("invalid value of field `since`")
	}

	// validate until
	if !p.Until.IsValid {
		return errors.New("invalid value of field `until`")
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
	Since *datatypes.ISO8601Date `json:"since"`

	// optional options
	Until *datatypes.ISO8601Date `json:"until"`
}

func (p *transactionSummaryRequest) Before() error {
	if p.Since == nil {
		return errors.New(
			schema.MissingRequiredFieldErrorDetail("since"),
		)
	}

	if p.Until == nil {
		p.Until = &datatypes.ISO8601Date{Time: time.Now()}
	}

	if !p.Since.IsValid {
		return errors.New("invalid value of field `since`")
	}

	if !p.Until.IsValid {
		return errors.New("invalid value of field `until`")
	}

	return nil
}

type transactionSummaryResponse struct {
	Count   int     `json:"count"`
	Income  float64 `json:"income"`
	Outcome float64 `json:"outcome"`
	Total   float64 `json:"total"`
}

// TransactionSummary returns summary for transactions owned by
// current user for given period in base currency (count, income, outcome, total).
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
		if transaction.Amount > 0 {
			income += transaction.BaseAmount
		} else {
			outcome += transaction.BaseAmount
		}
		total += transaction.BaseAmount
	}

	response := transactionSummaryResponse{
		Count:   len(transactions),
		Income:  income,
		Outcome: outcome,
		Total:   total,
	}
	request.SendSuccessfulResponse(&response)
}
