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
	"strconv"
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
	fmt.Println("benchmark: ", timeEnd-timeBegin)

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

func Test_ServiceBatch(t *testing.T) {
	clearDb()

	db1 := db1.NewDatabase(db1.NewConnect())
	db2 := db2.NewDatabase(db2.NewConnect())
	dbLog := dbLog.NewDatabase(dbLog.NewConnect())
	srv := service.NewService(db1, db2, dbLog)

	for i := 1; i <= 20; i++ {
		err := srv.CreateWallet(strconv.Itoa(i))
		if err != nil {
			log.Println(err)
		}
	}

	for i := 1; i <= 20; i++ {
		err := srv.UpdateBalance(strconv.Itoa(i), 3000)
		if err != nil {
			log.Println(err)
		}
	}

	walletSender := []string{"1", "2", "2", "2", "4", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "6"}
	walletRecipient := []string{"2", "1", "3", "4", "2", "9", "20", "19", "18", "17", "16", "15", "14", "13", "12", "10", "9", "8", "1", "20"}
	amount := []int{80, 30, 40, 50, 60, 70, 80, 90, 80, 70, 60, 50, 40, 30, 20, 10, 40, 70, 90, 80}

	timeBegin := time.Now().UnixMilli()
	srv.TransferBatch(walletSender, walletRecipient, amount)
	timeEnd := time.Now().UnixMilli()
	fmt.Println("benchmark: ", timeEnd-timeBegin)
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
