package main

import (
	"github.com/obscuronet/obscuro-playground/tools/transactioninjector"
	"os"
)

func main() {
	config := transactioninjector.ParseCLIArgs()
	transactioninjector.InjectTransactions(config)
	os.Exit(0)
}
