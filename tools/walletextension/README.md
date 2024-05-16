# Ten Gateway

Ten Gateway is available as a hosted service
(read more about it [here](https://docs.ten.xyz/docs/tools-infrastructure/hosted-gateway)),
but we also provide the source code for running it locally.
In this document, we will explain how to run and use Ten Gateway locally.

## Running Ten Gateway locally

Ten Gateway consists of two main parts:
- backend (which handles all api and rpc requests and is responsible for communication with the Ten node and encrypting/decrypting traffic)
- frontend (which is a web application that allows users to interact with the Ten network and manage their accounts in a user-friendly way through web browsers)

### Running the backend
To run Ten Gateway backend, you need to build and run the `walletextension` tool.

The precompiled binaries for macOS ARM64, macOS AMD64, Windows AMD64 and Linux AMD64 can be built by running the 
following commands from the `tools/walletextension/main` folder:

```
env GOOS=darwin GOARCH=amd64 go build -o ../bin/wallet_extension_macos_amd64 .
env GOOS=darwin GOARCH=arm64 go build -o ../bin/wallet_extension_macos_arm64 .
env GOOS=windows GOARCH=amd64 go build -o ../bin/wallet_extension_win_amd64.exe .
env GOOS=linux GOARCH=amd64 go build -o ../bin/wallet_extension_linux_amd64 .
```

The binaries will be created in the `tools/walletextension/bin` folder.
To run it you need to run built binary with the following command (depending on your OS):

```
./wallet_extension_macos_amd64

```

You can see the available flags (and default values) by running the binary with the `-h` flag.

#### Verifying the API

To verify that the API is working correctly, you can use `curl` to send a request to the `network-health` endpoint. This endpoint is accessible at `http://127.0.0.1:3000/v1/network-health`.

Open a terminal and run the following `curl` command:

```sh
curl -X GET http://127.0.0.1:3000/v1/network-health/
```

The Expected response is:
```json
{"id":"1","jsonrpc":"2.0","result":{"Errors":[],"OverallHealth":true}
```

### Running the frontend
To run Ten Gateway frontend you need to go to `go-ten/tools/walletextension/frontend` and run the following commands:

`npm run build` - to build the frontend (if you are using a different port/host you need to set environment variable `NEXT_PUBLIC_API_GATEWAY_URL` to match that change)

`npm run start` - to start the frontend



### HTTP Endpoints

For interacting with Ten Gateway, there are the following HTTP endpoints available:

- `GET /v1/join`

It generates and returns userID which needs to be added as a query parameter "u" to the URL in your Metamask
(or another provider) as it identifies you.

- `POST /v1/authenticate?token=$EncryptionToken`

With this endpoint, you submit a signed message in the format `Register <userID> for <account>`
from that account which proves that you hold private keys for it, and it links that account with your userID.

- `GET /v1/query/address?token=$EncryptionToken&a=$Address`

This endpoint responds with a JSON of true or false if the address "a" is already registered for user "u"


- `POST "/v1/revoke?token=$EncryptionToken"`

When this endpoint is triggered, the userId with the authenticated viewing keys should be deleted.
