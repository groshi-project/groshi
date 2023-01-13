package handles

import (
	"fmt"
	"github.com/jieggii/groshi/internal/database"
	"github.com/jieggii/groshi/internal/ghttp"
	"github.com/jieggii/groshi/internal/handles/schema"
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
	fmt.Println(params)
	transaction := database.Transaction{ // todo
		Amount:      params.Amount,
		Currency:    params.Currency,
		Description: params.Description,

		UUID:    "generated",
		OwnerId: currentUser.ID,
	}
	_, err := database.Db.NewInsert().Model(&transaction).Exec(database.Ctx)
	if err != nil {
		request.SendErrorResponse(schema.ServerSideError, "Could not create new transaction.", err)
		return
	}
	response := transactionCreateResponse{}
	request.SendSuccessResponse(&response)
}

type transactionReadRequest struct {
	UUID string `json:"uuid"`
}

type transactionReadResponse struct {
}

func TransactionRead(request *ghttp.Request, currentUser *database.User) {

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

}
