package database

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Context = context.TODO()

var Client *mongo.Client

var UsersCol *mongo.Collection
var TransactionsCol *mongo.Collection
var CurrencyRatesCol *mongo.Collection

func InitDatabase(host string, port int, username string, password string, databaseName string) error {
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%v:%v@%v:%v/", username, password, host, port),
	)
	var err error
	Client, err = mongo.Connect(Context, clientOptions)
	if err != nil {
		return err
	}
	err = Client.Ping(Context, nil)
	if err != nil {
		return err
	}
	database := Client.Database(databaseName)

	UsersCol = database.Collection("users")
	TransactionsCol = database.Collection("transactions")
	CurrencyRatesCol = database.Collection("currency-rates")

	return nil
}

// User represents service user.
type User struct {
	ID primitive.ObjectID `bson:"_id"`
	//UUID string             `bson:"uuid"`

	Username string `bson:"username"`
	Password string `bson:"password"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (u *User) JSON() gin.H {
	return gin.H{
		"username": u.Username,
	}
}

// Transaction represents financial transaction created by User.
type Transaction struct {
	ID   primitive.ObjectID `bson:"_id"`
	UUID string             `bson:"uuid"`

	OwnerID primitive.ObjectID `bson:"owner_id"`

	Amount      int       `bson:"amount"`   // amount of transaction in MINOR units
	Currency    string    `bson:"currency"` // currency code in ISO-4217 format
	Description string    `bson:"description"`
	Date        time.Time `bson:"date"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (t *Transaction) JSON() gin.H {
	return gin.H{
		"uuid": t.UUID,

		"amount":      t.Amount,
		"currency":    t.Currency,
		"description": t.Description,
		"date":        t.Date,

		"created_at": t.CreatedAt,
		"updated_at": t.UpdatedAt,
	}
}

// CurrencyRates TODO
type CurrencyRates struct {
	ID primitive.ObjectID `bson:"_id"`

	BaseCurrency string                 `bson:"currency"`
	Rates        map[string]interface{} `bson:"rate"`

	UpdatedAt time.Time `bson:"updated_at"`
}
