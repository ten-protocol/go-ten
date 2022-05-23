package main

import (
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
	"time"
)

func main() {
	l1Client, err := ethclient.NewEthClient()
	if err != nil {
		panic(err)
	}
	l2Client, err := obscuroclient.NewClient(0)

	txInjector := simulation.NewTransactionInjector(
		3,
		time.Second,
		nil,
		[]ethclient.EthClient{l1Client},
		[]*obscuroclient.Client{&l2Client},
	)

	txInjector.Start()
}
