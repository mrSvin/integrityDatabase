package dbLog

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var uriMongo = "mongodb://localhost:27017"
var dbName = "walletDb"
var dbCollection = "wallet_logs"

type logWallet struct {
	Id            string
	TypeOperation string
	OldBalance    int
	NewBalance    int
	TimeOperation int64
	Hash          string
}

func CreateLog(walletId string, oldBalance int, newBalance int, timeOperation int64, hash string, typeOperation string) error {
	clientOptions := options.Client().ApplyURI(uriMongo)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	collection := client.Database(dbName).Collection(dbCollection)

	logWallet := logWallet{
		Id:            walletId,
		TypeOperation: typeOperation,
		OldBalance:    oldBalance,
		NewBalance:    newBalance,
		TimeOperation: timeOperation,
		Hash:          hash,
	}

	_, err = collection.InsertOne(context.Background(), logWallet)
	if err != nil {
		return err
	}

	log.Println(time.Now(), " ", dbCollection, " wallet ", walletId, " created")
}
