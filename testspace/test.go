package main

import (
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	hostAndPort := "http://127.0.0.1:12000"

	//client, err := ethclient.Dial(hostAndPort)
	//if err != nil {
	//	panic(err)
	//}
	//
	//chainID, err := client.ChainID(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//
	//if chainID.Uint64() != 1337 {
	//	panic("Did not retrieve correct chain ID.")
	//}

	rpcClient, err := rpc.Dial(hostAndPort)
	if err != nil {
		panic(err)
	}

	var result string
	err = rpcClient.Call(&result, "obscuro_sendTransactionEncrypted")
	if err != nil {
		panic(err)
	}

	println(result)
}
