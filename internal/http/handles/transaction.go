package handles

import (
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/http/ghttp"
	"github.com/jieggii/groshi/internal/http/schema"
	"time"
)

type transactionCreateRequest struct {
	Amount      float64           `json:"amount,string"` // todo: think about string and int
	Currency    database.Currency `json:"currency"`
	Description string            `json:"description"`
}
type transactionCreateResponse struct {
	UUID string `json:"uuid"`
}

func TransactionCreate(request *ghttp.Request, currentUser *database.User) {
	params := transactionCreateRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	transaction := database.Transaction{
		Amount:      params.Amount,
		Currency:    params.Currency,
		Description: params.Description,

		OwnerId: currentUser.ID,
	}
	_, err := database.Db.NewInsert().Model(&transaction).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse(
			"Could not create new transaction.", err,
		)
		return
	}

	response := transactionCreateResponse{UUID: transaction.UUID}
	request.SendSuccessResponse(&response)
}

type transactionReadRequest struct {
	UUID string `json:"uuid"`
}

type transactionReadResponse struct {
	UUID string `json:"uuid"`

	Amount      float64           `json:"amount"`
	Currency    database.Currency `json:"currency"`
	Description string            `json:"description"`

	Owner string    `json:"owner"`
	Date  time.Time `json:"date"`
}

func TransactionRead(request *ghttp.Request, currentUser *database.User) {
	params := transactionReadRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	transaction, err := database.FetchTransactionByUUID(params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(schema.TransactionNotFound)
		return
	}

	if transaction.Owner.ID != currentUser.ID && !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(schema.AccessDenied)
		return
	}

	response := transactionReadResponse{
		UUID:        transaction.UUID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,

		Owner: transaction.Owner.Username,
		Date:  transaction.Date,
	}
	request.SendSuccessResponse(&response)
}

type transactionUpdateRequest struct {
}

type transactionUpdateResponse struct {
}

func TransactionUpdate(request *ghttp.Request, currentUser *database.User) {

}

type transactionDeleteRequest struct {
	UUID string `json:"uuid"`
}

type transactionDeleteResponse struct {
}

func TransactionDelete(request *ghttp.Request, currentUser *database.User) {
	params := transactionDeleteRequest{}
	if ok := request.DecodeSafe(&params); !ok {
		return
	}

	transaction, err := database.FetchTransactionByUUID(params.UUID)
	if err != nil {
		request.SendClientSideErrorResponse(schema.TransactionNotFound)
		return
	}
	if currentUser.ID != transaction.OwnerId && !currentUser.IsSuperuser {
		request.SendClientSideErrorResponse(schema.AccessDenied)
		return
	}
	_, err = database.Db.NewDelete().Model(&transaction).Where("uuid = ", transaction.UUID).Exec(database.Ctx)
	if err != nil {
		request.SendServerSideErrorResponse("Could not delete transaction.", err)
		return
	}
	response := transactionDeleteResponse{}
	request.SendSuccessResponse(&response)
}
