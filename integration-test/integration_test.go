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
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var uriMongo = "mongodb://localhost:27017"
var countWallets = 100 //for batch test

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

	for i := 1; i <= countWallets; i++ {
		err := srv.CreateWallet(strconv.Itoa(i))
		if err != nil {
			log.Println(err)
		}
	}

	for i := 1; i <= countWallets; i++ {
		err := srv.UpdateBalance(strconv.Itoa(i), 5000)
		if err != nil {
			log.Println(err)
		}
	}

	walletSender, walletRecipient, amount := generateRandomTransfer()

	timeBegin := time.Now().UnixMilli()
	//кидаем 3 пула с группой трансферов
	srv.TransferBatch(walletSender, walletRecipient, amount)
	srv.TransferBatch(walletSender, walletRecipient, amount)
	srv.TransferBatch(walletSender, walletRecipient, amount)
	timeEnd := time.Now().UnixMilli()
	countTime := 1000 / float32(timeEnd-timeBegin)
	benchmark := countTime * float32(countWallets*3)
	fmt.Println("benchmark: ", benchmark, " tps")
}

func clearDb() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uriMongo))
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
	log.Println("база очищена")
}

func generateRandomTransfer() ([]string, []string, []int) {
	var walletSender []string
	var walletRecipient []string
	var amount []int

	for i := 0; i < countWallets; i++ {
		walletSender = append(walletSender, strconv.Itoa(rand.Intn(countWallets)+1))
		walletRecipient = append(walletRecipient, strconv.Itoa(rand.Intn(countWallets)+1))
		amount = append(amount, rand.Intn(40)+10)
	}

	return walletSender, walletRecipient, amount
}
