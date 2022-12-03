package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var users *mongo.Collection
var transactions *mongo.Collection

var ctx = context.TODO()

func buildURI(host string, port int) string {
	return fmt.Sprintf("mongodb://%v:%v/", host, port)
}

func Connect(host string, port int, dbName string) error {
	clientOptions := options.Client().ApplyURI(buildURI(host, port))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	db := client.Database(dbName)
	users = db.Collection("users")
	transactions = db.Collection("transactions")
	return nil
}
