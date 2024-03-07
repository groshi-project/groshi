package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/groshi-project/groshi/internal/database"
	"net/http"
	"time"
)

type transactionsCreateParams struct {
	Amount      int32     `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type transactionsCreateResponse struct {
	UUID string `json:"uuid"`
}

func (s *Service) TransactionsCreate(w http.ResponseWriter, r *http.Request) {
	params := transactionsCreateParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// fetch provided currency:
	currency := database.Currency{}
	if err := s.Database.Client.NewSelect().Model(database.EmptyCurrency).Where("code = ?", params.Currency).Scan(s.Database.Ctx, &currency); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fetch current user: (todo: fetch only user's id)
	claims, err := s.JWTAuthority.ExtractClaims(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := claims["username"].(string)
	user := database.User{}
	if err := s.Database.Client.NewSelect().Model(&user).Where("username = ?", username).Scan(s.Database.Ctx, &user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check timestamp timezone:
	location := params.Timestamp.Location().String()
	_, err = time.LoadLocation(location)
	if err != nil || location == "" || location == "Local" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// create a new transaction owned by the current user:
	transaction := database.Transaction{
		Amount:      params.Amount,
		CurrencyID:  currency.ID,
		Description: params.Description,
		OwnerID:     user.ID,
		Timestamp:   params.Timestamp,
		Timezone:    location,
	}
	if _, err := s.Database.Client.NewInsert().Model(&transaction).Exec(s.Database.Ctx); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// respond:
	response := transactionsCreateResponse{
		UUID: transaction.UUID.String(),
	}
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Service) TransactionsGetOne(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) TransactionsGet(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) TransactionsUpdate(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) TransactionsDelete(w http.ResponseWriter, r *http.Request) {

}
