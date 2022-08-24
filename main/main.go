package main

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	url := "ws://127.0.0.1:37500/"
	client, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ch := make(chan *string)
	sub, err := client.Subscribe(context.Background(), "eth", ch, "aSubAPI", nil)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-ch:
			println("jjj received message:", *msg)
		case err := <-sub.Err():
			panic(err)
		}
	}
}
