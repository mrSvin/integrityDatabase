package main

import (
	"fmt"
	"integrity/service"
	"log"
)

func main() {

	service.CreateWallet("1")
	balance, err := service.ReadWalletBalance("1")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(balance)

}
