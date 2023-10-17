package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"integrity/db1"
	"integrity/db2"
	"integrity/dbLog"
	"integrity/internal"
	"time"
)

type service struct {
	db1   *db1.Database
	db2   *db2.Database
	dbLog *dbLog.Database
}

func NewService(db1 *db1.Database, db2 *db2.Database, dbLog *dbLog.Database) *service {
	return &service{
		db1:   db1,
		db2:   db2,
		dbLog: dbLog,
	}
}

func (s *service) CreateWallet(walletId string) error {

	time := time.Now().UnixNano()
	hashString := getHash(walletId, 0, time)
	hashLength := len(hashString) / 2
	hashBegin := hashString[:hashLength]
	hashEnd := hashString[hashLength:]

	err := s.db1.CreateWallet(walletId, time, hashBegin)
	if err != nil {
		return err
	}
	err = s.db2.CreateWallet(walletId, hashEnd)
	if err != nil {
		return err
	}
	err = s.dbLog.CreateLog(walletId, 0, 0, time, hashBegin, "create")
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ReadWalletBalance(walletId string) (int, error) {
	err := s.checkWalletHash(walletId)
	if err != nil {
		return 0, err
	}
	data, err := s.db1.ReadWallet(walletId)
	return data.Balance, err
}

func (s *service) ReadWalletBalanceForTransfer(walletId string) (int, error) {
	data, err := s.db1.ReadWallet(walletId)
	return data.Balance, err
}

// Для эмиссии, уничтожения валюты
func (s *service) UpdateBalance(walletId string, newBalance int) error {
	err := s.checkWalletHash(walletId)
	if err != nil {
		return err
	}

	timeUpdate := time.Now().UnixNano()
	hashString := getHash(walletId, newBalance, timeUpdate)
	hashLength := len(hashString) / 2
	hashBegin := hashString[:hashLength]
	hashEnd := hashString[hashLength:]

	walletInfo, err := s.db1.ReadWallet(walletId)
	if err != nil {
		return err
	}

	err = s.db1.UpdateBalanceWallet(walletId, newBalance, timeUpdate, hashBegin)
	if err != nil {
		return err
	}
	err = s.db2.UpdateHashWallet(walletId, hashEnd)
	if err != nil {
		return err
	}
	err = s.dbLog.CreateLog(walletId, walletInfo.Balance, newBalance, timeUpdate, hashBegin, "update")
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Transfer(walletIdSender string, walletIdRecipient string, sendMoney int) error {

	senderBalanceCh := make(chan int)
	recipientBalanceCh := make(chan int)

	// запускаем две горутины для чтения баланса отправителя и получателя
	go func() {
		balanceSender, err := s.ReadWalletBalanceForTransfer(walletIdSender)
		if err != nil {
			// отправляем ошибку в канал
			senderBalanceCh <- -1
			return
		}
		senderBalanceCh <- balanceSender
	}()

	go func() {
		balanceRecipient, err := s.ReadWalletBalanceForTransfer(walletIdRecipient)
		if err != nil {
			// отправляем ошибку в канал
			recipientBalanceCh <- -1
			return
		}
		recipientBalanceCh <- balanceRecipient
	}()

	// ожидаем результаты чтения баланса отправителя и получателя
	balanceSender := <-senderBalanceCh
	balanceRecipient := <-recipientBalanceCh

	// проверяем наличие ошибок при чтении баланса
	if balanceSender == -1 || balanceRecipient == -1 {
		return errors.New("failed to read wallet balance")
	}

	// запускаем две горутины для обновления баланса отправителя и получателя
	errCh := make(chan error)
	go func() {
		err := s.UpdateBalance(walletIdSender, balanceSender-sendMoney)
		errCh <- err
	}()

	go func() {
		err := s.UpdateBalance(walletIdRecipient, balanceRecipient+sendMoney)
		errCh <- err
	}()

	// ожидаем результаты обновления баланса отправителя и получателя
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			return err
		}
	}

	return nil
}

func (s *service) TransferBatch(walletIdSenderArray []string, walletIdRecipientArray []string, amount []int) {
	newWalletSender, newWalletRecipient, newAmount := internal.CalculateMirrorTransfer(walletIdSenderArray, walletIdRecipientArray, amount)

	fmt.Println(newWalletSender)
	fmt.Println(newWalletRecipient)
	fmt.Println(newAmount)

	//TODO add async run
	for i := 0; i < len(newWalletSender); i++ {
		s.Transfer(newWalletSender[i], newWalletRecipient[i], newAmount[i])
	}
}

func (s *service) checkWalletHash(walletId string) error {

	walletDb1, err := s.db1.ReadWallet(walletId)
	if err != nil {
		return err
	}
	walletDb2, err := s.db2.ReadWallet(walletId)
	if err != nil {
		return err
	}

	hashNode := walletDb1.Hash + walletDb2.Hash
	hashDb := getHash(walletDb1.Id, walletDb1.Balance, walletDb1.TimeOperation)

	if hashNode != hashDb {
		return errors.New("invalid hash in wallet id: " + walletId)
	}

	return nil
}

func getHash(walletId string, balance int, time int64) string {
	dataHash := fmt.Sprintf("%s %v %v", walletId, balance, time)
	hash := sha256.Sum256([]byte(dataHash))
	return fmt.Sprintf("%x\n", hash)
}
