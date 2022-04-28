package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// RPC request handler.
	// This is a helpful resource: https://docs.alchemy.com/alchemy/apis/ethereum/eth_chainid.
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		var jsonMap map[string]interface{}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			fmt.Println(err)
		}

		switch jsonMap["method"] {
		case "eth_blockNumber":
			fmt.Println("Received eth_blockNumber request.")
			_, err = resp.Write([]byte("{\"result\": \"999\"}"))
			if err != nil {
				fmt.Println(err)
			}
		case "eth_chainId":
			fmt.Println("Received eth_chainId request.")
			_, err = resp.Write([]byte("{\"result\": \"0x309\"}"))
			if err != nil {
				fmt.Println(err)
			}
		case "eth_gasPrice":
			fmt.Println("Received eth_gasPrice request.")
			_, err = resp.Write([]byte("{\"result\": \"0\"}"))
			if err != nil {
				fmt.Println(err)
			}
		case "eth_estimateGas":
			fmt.Println("Received eth_estimateGas request.")
			_, err = resp.Write([]byte("{\"result\": \"0\"}"))
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Print("Received unhandled request:")
			fmt.Println(jsonMap["method"])
		}
	})

	// Web app handler.
	http.Handle("/register/", http.StripPrefix("/register/", http.FileServer(http.Dir("./tools/clientproxy/main/static"))))
	http.HandleFunc("/generateViewingKeyPair", func(resp http.ResponseWriter, req *http.Request) {
		// todo - generate viewing key properly
		// todo - store private key
		_, err := resp.Write([]byte("dummyViewingKey"))
		if err != nil {
			fmt.Println(err)
		}
	})
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
