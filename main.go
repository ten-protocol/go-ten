package main

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
)

const url = "ws://127.0.0.1:37500/" // The websocket address for the first Obscuro node in the full network simulation.

func main() {
	client, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ch := make(chan *string)
	var subArgs []interface{} // By passing no additional args, we subscribe specifically for newly-mined blocks.
	sub, err := client.Subscribe(context.Background(), "eth", ch, "logs", subArgs)
	if err != nil {
		panic(err)
	}

	println("jjj subscribed")

	for {
		select {
		case msg := <-ch:
			println("jjj received message:", *msg)
		case err := <-sub.Err():
			panic(err)
		}
	}
}
