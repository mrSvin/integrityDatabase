package integration_test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"integrity/db1"
	"integrity/db2"
	"integrity/dbLog"
	"integrity/service"
	"log"
	"testing"
	"time"
)

func Test_Service(t *testing.T) {

	clearDb()

	db1 := db1.NewDatabase(db1.NewConnect())
	db2 := db2.NewDatabase(db2.NewConnect())
	dbLog := dbLog.NewDatabase(dbLog.NewConnect())
	srv := service.NewService(db1, db2, dbLog)

	err := srv.CreateWallet("1")
	if err != nil {
		log.Println(err)
	}

	err = srv.CreateWallet("2")
	if err != nil {
		log.Println(err)
	}

	balance, err := srv.ReadWalletBalance("1")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(balance)

	err = srv.UpdateBalance("1", 10000)
	if err != nil {
		log.Println(err)
	}

	balance, err = srv.ReadWalletBalance("1")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(balance)

	timeBegin := time.Now().UnixMilli()
	for i := 0; i < 5000; i++ {
		srv.Transfer("1", "2", 1)
		srv.Transfer("2", "1", 1)
	}
	timeEnd := time.Now().UnixMilli()
	fmt.Println(timeEnd - timeBegin)

	balance, err = srv.ReadWalletBalance("2")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("balance wallet id2: ", balance)

	balance, err = srv.ReadWalletBalance("1")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("balance wallet id1: ", balance)

}

func clearDb() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Удаление коллекций
	err = client.Database("walletDb").Collection("wallet_logs").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database("walletDb").Collection("wallet_node_1").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database("walletDb").Collection("wallet_node_2").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
