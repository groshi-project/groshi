package handles

import (
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/ghttp/schema"
	"time"
)

type transactionCreateRequest struct {
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (p *transactionCreateRequest) validate() bool {
	return p.Currency != "" && p.Amount >= 0
}

type transactionCreateResponse struct {
	UUID string `json:"uuid"`
}

func TransactionCreate(request *ghttp.Request, currentUser *database.User) {
	params := transactionCreateRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
		)
		return
	}

	transaction := database.Transaction{
		Amount:      params.Amount,
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
	request.SendSuccessResponse(&response)
}

type transactionReadRequest struct {
	UUID string `json:"uuid"`
}

func (p *transactionReadRequest) validate() bool {
	return p.UUID != ""
}

type transactionReadResponse struct {
	UUID string `json:"uuid"`

	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`

	Owner string    `json:"owner"`
	Date  time.Time `json:"date"`
}

func TransactionRead(request *ghttp.Request, currentUser *database.User) {
	params := transactionReadRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
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
			schema.AccessDeniedErrorTag, schema.NoRightToPerformOperationErrorDetail,
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
	}
	request.SendSuccessResponse(&response)
}

type transactionUpdateRequest struct {
}

func (p *transactionUpdateRequest) validate() bool {
	return true
}

type transactionUpdateResponse struct {
}

func TransactionUpdate(request *ghttp.Request, currentUser *database.User) {

}

type transactionDeleteRequest struct {
	UUID string `json:"uuid"`
}

func (p *transactionDeleteRequest) validate() bool {
	return p.UUID != ""
}

type transactionDeleteResponse struct {
	UUID string `json:"uuid"`
}

func TransactionDelete(request *ghttp.Request, currentUser *database.User) {
	params := transactionDeleteRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	if !params.validate() {
		request.SendClientSideErrorResponse(
			schema.InvalidRequestErrorTag, schema.RequestBodyDidNotPassValidation,
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
			schema.AccessDeniedErrorTag, schema.NoRightToPerformOperationErrorDetail,
		)
		return
	}

	_, err = database.Db.NewDelete().Model(transaction).Where("uuid = ?", transaction.UUID).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("could not delete transaction", err)
		return
	}

	response := transactionDeleteResponse{UUID: transaction.UUID}
	request.SendSuccessResponse(&response)
}
