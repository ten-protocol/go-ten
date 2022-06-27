package main

import (
	"os"

	"github.com/obscuronet/obscuro-playground/tools/transactioninjector"
)

func main() {
	config := transactioninjector.ParseCLIArgs()
	transactioninjector.InjectTransactions(config)
	os.Exit(0)
}
