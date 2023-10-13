package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"integrity/db1"
	"integrity/db2"
	"integrity/dbLog"
	"log"
	"time"
)

func CreateWallet(walletId string) error {

	time := time.Now().UnixNano()
	hashString := getHash(walletId, 0, time)
	hashLength := len(hashString) / 2
	hashBegin := hashString[:hashLength]
	hashEnd := hashString[hashLength:]

	err := db1.CreateWallet(walletId, time, hashBegin)
	if err != nil {
		return err
	}
	err = db2.CreateWallet(walletId, hashEnd)
	if err != nil {
		return err
	}
	err = dbLog.CreateLog(walletId, 0, 0, time, hashBegin, "create")
	if err != nil {
		return err
	}
	return nil
}

func ReadWalletBalance(walletId string) (int, error) {
	err := checkWalletHash(walletId)
	if err != nil {
		return 0, err
	}
	data, err := db1.ReadWallet(walletId)
	return data.Balance, err
}

// Для эмиссии и уничтожения валюты
func UpdateBalance(walletId string, newBalance int) error {
	err := checkWalletHash(walletId)
	if err != nil {
		return err
	}

	time := time.Now().UnixNano()
	hashString := getHash(walletId, 0, time)
	hashLength := len(hashString) / 2
	hashBegin := hashString[:hashLength]
	hashEnd := hashString[hashLength:]

	walletInfo, err := db1.ReadWallet(walletId)
	if err != nil {
		return err
	}

	err = db1.UpdateBalanceWallet(walletId, newBalance, time, hashBegin)
	if err != nil {
		return err
	}
	err = db2.UpdateHashWallet(walletId, hashEnd)
	if err != nil {
		return err
	}
	err = dbLog.CreateLog(walletId, walletInfo.Balance, newBalance, time, hashBegin, "update")
	if err != nil {
		return err
	}

	return nil
}

func checkWalletHash(walletId string) error {

	walletDb1, err := db1.ReadWallet(walletId)
	if err != nil {
		log.Println(err)
	}

	walletDb2, err := db2.ReadWallet(walletId)
	if err != nil {
		log.Println(err)
	}

	hashNode := walletDb1.Hash + walletDb2.Hash
	hashDb := getHash(walletDb1.Id, walletDb1.Balance, walletDb1.TimeOperation)

	if hashNode != hashDb {
		return errors.New("invalid hash " + walletId)
	}

	return nil
}

func getHash(walletId string, balance int, time int64) string {
	dataHash := fmt.Sprintf("%s %v %v", walletId, balance, time)
	hash := sha256.Sum256([]byte(dataHash))
	return fmt.Sprintf("%x\n", hash)
}
