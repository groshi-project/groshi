package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Context = context.TODO()

var Client *mongo.Client

// UsersCol is `users` collection.
var UsersCol *mongo.Collection

// TransactionsCol is `transactions` collection.
var TransactionsCol *mongo.Collection

// CurrencyRatesCol is `currency-rates` collection.
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
