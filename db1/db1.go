package db1

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var uriMongo = "mongodb://localhost:27017"
var dbName = "walletDb"
var dbCollection = "wallet_node_1"

type Wallet struct {
	Id            string
	Balance       int
	TimeOperation int64
	Hash          string
}

func CreateWallet(walletId string, timeCreate int64, hash string) error {
	clientOptions := options.Client().ApplyURI(uriMongo)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection(dbCollection)

	wallet := Wallet{
		Id:            walletId,
		Balance:       0,
		TimeOperation: timeCreate,
		Hash:          hash,
	}

	err = collection.FindOne(context.Background(), bson.M{"id": walletId}).Decode(&wallet)
	if err == nil {
		log.Println(time.Now(), " ", dbCollection, " wallet ", walletId, " already exists")
		return errors.New("db1 already exists wallet " + walletId)
	}

	_, err = collection.InsertOne(context.Background(), wallet)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(time.Now(), " ", dbCollection, " wallet ", walletId, " created")
	return nil
}

func ReadWallet(walletId string) (*Wallet, error) {
	clientOptions := options.Client().ApplyURI(uriMongo)

	// установка соединения с базой данных
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// проверка соединения с базой данных
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// выбор коллекции для чтения данных
	collection := client.Database(dbName).Collection(dbCollection)

	var wallet Wallet
	err = collection.FindOne(context.Background(), bson.M{"id": walletId}).Decode(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func UpdateBalanceWallet(walletId string, newBalance int, timeUpdate int64, hash string) error {
	clientOptions := options.Client().ApplyURI(uriMongo)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection(dbCollection)

	filter := bson.M{"id": walletId}
	update := bson.M{"$set": bson.M{"balance": newBalance, "timeOperation": timeUpdate, "hash": hash}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	log.Println(time.Now(), " ", dbCollection, " wallet ", walletId, " updated")
	return nil
}
