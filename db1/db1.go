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

type Database struct {
	Collection *mongo.Collection
}

type Wallet struct {
	Id            string
	Balance       int
	TimeOperation int64
	Hash          string
}

func (db *Database) CreateWallet(walletId string, timeCreate int64, hash string) error {

	wallet := Wallet{
		Id:            walletId,
		Balance:       0,
		TimeOperation: timeCreate,
		Hash:          hash,
	}

	err := db.Collection.FindOne(context.Background(), bson.M{"id": walletId}).Decode(&wallet)
	if err == nil {
		log.Println(time.Now(), " ", dbCollection, " wallet ", walletId, " already exists")
		return errors.New("db1 already exists wallet " + walletId)
	}

	_, err = db.Collection.InsertOne(context.Background(), wallet)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(time.Now(), " ", dbCollection, " wallet ", walletId, " created")
	return nil
}

func (db *Database) ReadWallet(walletId string) (*Wallet, error) {

	var wallet Wallet
	err := db.Collection.FindOne(context.Background(), bson.M{"id": walletId}).Decode(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (db *Database) UpdateBalanceWallet(walletId string, newBalance int, timeUpdate int64, hash string) error {
	filter := bson.M{"id": walletId}
	update := bson.M{"$set": bson.M{"balance": newBalance, "timeOperation": timeUpdate, "hash": hash}}

	_, err := db.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	log.Println(time.Now(), " ", dbCollection, " wallet ", walletId, " updated")
	return nil
}

func (db *Database) UpdateBatchBalanceWallet(walletId []string, newBalance []int, timeUpdate int64, hash []string) error {

	var bulkOps []mongo.WriteModel

	for i := 0; i < len(walletId); i++ {
		filter := bson.M{"id": walletId[i]}
		update := bson.M{"$set": bson.M{"balance": newBalance[i], "timeOperation": timeUpdate, "hash": hash[i]}}
		updateOne := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		bulkOps = append(bulkOps, updateOne)
	}
	_, err := db.Collection.BulkWrite(context.Background(), bulkOps)
	if err != nil {
		return err
	}
	return nil
}

func NewConnect() *mongo.Collection {
	clientOptions := options.Client().ApplyURI(uriMongo)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(dbName).Collection(dbCollection)
}

func NewDatabase(collection *mongo.Collection) *Database {
	return &Database{
		Collection: collection,
	}
}
