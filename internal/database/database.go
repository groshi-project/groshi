package database

import (
	"context"
	"fmt"
	"github.com/jieggii/groshi/internal/loggers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Context = context.TODO()

var Users *mongo.Collection
var Transactions *mongo.Collection

func InitDatabase(host string, port int, username string, password string, databaseName string) error {
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%v:%v@%v:%v/", username, password, host, port),
	)
	client, err := mongo.Connect(Context, clientOptions)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Disconnect(Context)
		if err != nil {
			loggers.Error.Fatal(err)
		}
	}()

	err = client.Ping(Context, nil)
	if err != nil {
		return err
	}
	database := client.Database(databaseName)

	Users = database.Collection("users")
	Transactions = database.Collection("transactions")

	return nil
}

func GenerateUUID() string {
	return "123-123"
}

type User struct {
	ID   primitive.ObjectID `bson:"_id"`
	UUID string             `bson:"uuid"`

	Username string `bson:"username"`
	Password string `bson:"password"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Transaction struct {
	ID primitive.ObjectID `bson:"_id"`

	UUID      string `bson:"uuid"`
	OwnerUUID string `bson:"owner_uuid"`

	Amount      float64   `bson:"amount"`
	Currency    string    `bson:"currency"`
	Description string    `bson:"description"`
	Date        time.Time `bson:"date"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

//func FindOne(collection *mongo.Collection, query bson.D, v interface{}) found {
//	err := collection.FindOne(Context, query).Decode(&v)
//	return err
//}
//
//func Exists(collection *mongo.Collection, query bson.D) (bool, error) {
//	err := collection.FindOne(Context, query).Err()
//	if err != nil {
//		if errors.Is(err, mongo.ErrNoDocuments) {
//			return false, nil
//		}
//		return false, err
//	}
//	return true, nil
//}
//
//func InsertOne(collection *mongo.Collection, obj interface{}) {
//	collection.InsertOne(Context, obj)
//	//result, err := collection.InsertOne(
//	//	context.TODO(),
//	//	bson.D{
//	//		{"animal", "Dog"},
//	//		{"breed", "Beagle"}
//	//	}
//	//)
//}
