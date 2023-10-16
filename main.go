package main

import (
	"fmt"
	"integrity/db1"
	"integrity/db2"
	"integrity/dbLog"
	"integrity/service"
	"log"
)

func main() {
	db1 := db1.NewDatabase(db1.NewConnect())
	db2 := db2.NewDatabase(db2.NewConnect())
	dbLog := dbLog.NewDatabase(dbLog.NewConnect())
	srv := service.NewService(db1, db2, dbLog)

	err := srv.CreateWallet("1")
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
}

//fmt.Println(getHash("1", 10000, 1697200252701001700))

//func getHash(walletId string, balance int, time int64) string {
//	dataHash := fmt.Sprintf("%s %v %v", walletId, balance, time)
//	hash := sha256.Sum256([]byte(dataHash))
//	return fmt.Sprintf("%x\n", hash)
//}
