This package contains a fork of the rpc stack from go-ethereum.

We need to change it, to add a custom authentication mechanism.

The RPC URL to the TEN Gateway contains a "token" URL parameter kept secret by every user, which we use for authentication.

1. Create an http service: `extract_params_handler.go` which extracts it and adds it to the Context
2. From the Context, it can be read by the http rpc handler
3. For WS, we need to add it to a couple of more objects to make it available during processing.
4. Removed unnecessary elements from the "node", such as p2p, database, etc
5. Enable full duplex to enable keep-alive connections