package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Context = context.TODO()

var Client *mongo.Client

var Users *mongo.Collection
var Transactions *mongo.Collection

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

	Users = database.Collection("users")
	Transactions = database.Collection("transactions")

	return nil
}

// GenerateUUID generates new UUID v4.
func GenerateUUID() string {
	return uuid.New().String()
}

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	//UUID string             `bson:"uuid"`

	Username string `bson:"username"`
	Password string `bson:"password"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

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
