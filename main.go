package main

import (
	"fmt"
)

func main() {

	walletIdSenderArray := []string{"wallet1", "wallet2", "wallet3"}
	walletIdRecipientArray := []string{"wallet4", "wallet1", "wallet6"}
	amount := []int{100, 200, 300}

	balances := calculateUpdateWalletBalances(walletIdSenderArray, walletIdRecipientArray, amount)

	for walletId, updateBalance := range balances {
		fmt.Println(walletId, updateBalance)
	}
}

func calculateUpdateWalletBalances(walletIdSenderArray []string, walletIdRecipientArray []string, amount []int) map[string]int {
	balances := make(map[string]int)
	for i := 0; i < len(walletIdSenderArray); i++ {
		sender := walletIdSenderArray[i]
		recipient := walletIdRecipientArray[i]
		transferAmount := amount[i]
		balances[sender] -= transferAmount
		balances[recipient] += transferAmount
	}
	result := make(map[string]int)
	for walletId, balance := range balances {
		if balance != 0 {
			result[walletId] = balance
		}
	}
	return result
}
