package integration_test

import (
	"fmt"
	"integrity/db1"
	"integrity/db2"
	"integrity/dbLog"
	"integrity/service"
	"log"
	"testing"
	"time"
)

func Test_Service(t *testing.T) {
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
	for i := 0; i < 10000; i++ {
		srv.Transfer("1", "2", 1)
	}
	timeEnd := time.Now().UnixMilli()
	fmt.Println(timeEnd - timeBegin)
	balance, err = srv.ReadWalletBalance("2")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(balance)

}
