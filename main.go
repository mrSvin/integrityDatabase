package main

import (
	"fmt"
	"integrity/internal"
)

func main() {

	in1 := []string{"1", "2", "2", "2", "4", "7"}
	in2 := []string{"2", "1", "3", "4", "2", "8"}
	in3 := []int{80, 30, 40, 50, 60, 70}

	newWalletSender, newWalletRecipient, newAmount := internal.CalculateMirrorTransfer(in1, in2, in3)
	fmt.Println(newWalletSender)
	fmt.Println(newWalletRecipient)
	fmt.Println(newAmount)
}
