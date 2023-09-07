package database

import (
	"github.com/groshi-project/groshi/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User represents service user.
type User struct {
	ID primitive.ObjectID `bson:"_id"`

	Username string `bson:"username"`
	Password string `bson:"password"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (u *User) APIModel() *models.User {
	return &models.User{Username: u.Username}
}

// Transaction represents financial transaction created by User.
type Transaction struct {
	ID   primitive.ObjectID `bson:"_id"`
	UUID string             `bson:"uuid"`

	OwnerID primitive.ObjectID `bson:"owner_id"`

	Amount      int    `bson:"amount"`   // amount of transaction in MINOR units
	Currency    string `bson:"currency"` // currency code in ISO-4217 format
	Description string `bson:"description"`

	Time time.Time `bson:"time"` // transaction timestamp (when it happened)

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (t *Transaction) APIModel() *models.Transaction {
	// please note: all time.Time values are returned relative to the UTC timezone!
	return &models.Transaction{
		UUID: t.UUID,

		Amount:      t.Amount,
		Currency:    t.Currency,
		Description: t.Description,
		Time:        t.Time.In(time.UTC).Format(time.RFC3339),

		CreatedAt: t.CreatedAt.In(time.UTC).Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.In(time.UTC).Format(time.RFC3339),
	}
}

// CurrencyRates TODO
type CurrencyRates struct {
	ID primitive.ObjectID `bson:"_id"`

	BaseCurrency string                 `bson:"currency"`
	Rates        map[string]interface{} `bson:"rates"`

	UpdatedAt time.Time `bson:"updated_at"`
}
