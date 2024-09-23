
# TEN Gateway Documentation

For a comprehensive overview, refer to the [official documentation](https://docs.ten.xyz/docs/tools-infrastructure/hosted-gateway).

## Running the Gateway Locally

### Backend

To run the backend locally, it is recommended to use **port 1443** to avoid conflicts with the frontend service, which typically runs on port 3000. First, build the backend using the `go build` command. Navigate to the `tools/walletextension/main` folder and use the following commands to build for your respective operating system:

```bash
# macOS AMD64
env GOOS=darwin GOARCH=amd64 go build -o ../bin/wallet_extension_macos_amd64 .

# macOS ARM64
env GOOS=darwin GOARCH=arm64 go build -o ../bin/wallet_extension_macos_arm64 .

# Windows AMD64
env GOOS=windows GOARCH=amd64 go build -o ../bin/wallet_extension_win_amd64.exe .

# Linux AMD64
env GOOS=linux GOARCH=amd64 go build -o ../bin/wallet_extension_linux_amd64 .
```

The binaries will be available in the `tools/walletextension/bin` directory.
Run the compiled binary and specify the desired port.
Example:

```bash
./wallet_extension_macos_arm64 --port 1443
```

### Additional Backend Configuration Options

- **`--host`**: The host where the wallet extension should open the port. Default: `127.0.0.1`.
- **`--port`**: The port on which to serve the wallet extension. Default: `3000`.
- **`--portWS`**: The port on which to serve websocket JSON RPC requests. Default: `3001`.
- **`--nodeHost`**: The host on which to connect to the Obscuro node. Default: `erpc.sepolia-testnet.ten.xyz`.
- **`--nodePortHTTP`**: The port on which to connect to the Obscuro node via RPC over HTTP. Default: `80`.
- **`--nodePortWS`**: The port on which to connect to the Obscuro node via RPC over websockets. Default: `81`.
- **`--logPath`**: The path to use for the wallet extension's log file. Default: `sys_out`.
- **`--databasePath`**: The path for the wallet extension's database file. Default: `.obscuro/gateway_database.db`.
- **`--verbose`**: Flag to enable verbose logging of wallet extension traffic. Default: `false`.
- **`--dbType`**: Define the database type (`sqlite` or `mariaDB`). Default: `sqlite`.
- **`--dbConnectionURL`**: If `dbType` is set to `mariaDB`, this must be set.
- **`--tenChainID`**: ChainID of the TEN network that the gateway is communicating with. Default: `443`.
- **`--storeIncomingTxs`**: Flag to enable storing incoming transactions in the database for debugging purposes. Default: `true`.
- **`--rateLimitUserComputeTime`**: Represents how much compute time a user is allowed to use within the `rateLimitWindow` time. Set to `0` to disable rate limiting. Default: `10s`.
- **`--rateLimitWindow`**: Time window in which a user is allowed to use the defined compute time. Default: `1m`.
- **`--maxConcurrentRequestsPerUser`**: Number of concurrent requests allowed per user. Default: `3`.


### Frontend

Once the backend is running, navigate to the `tools/walletextension/frontend` directory and execute the following commands:

```bash
npm install
npm run dev
```

The frontend will be accessible on `http://localhost:3000`.

## HTTP Endpoints

TEN Gateway exposes several HTTP endpoints for interaction:

- **`GET /v1/join`**  
  Generates and returns a `userID`, which needs to be added as a query parameter `u` in your Metamask (or another provider) URL to identify you.

- **`POST /v1/authenticate?token=$EncryptionToken`**  
  Submits a signed message in the format `Register <userID> for <account>`, proving ownership of the private keys for the account, and links that account with the `userID`.

- **`GET /v1/query/address?token=$EncryptionToken&a=$Address`**  
  Returns a JSON response indicating whether the address "a" is registered for the user "u".

- **`POST /v1/revoke?token=$EncryptionToken`**  
  Deletes the userId along with the associated authenticated viewing keys.

- **`GET /v1/health`**  
  Returns a health status of the service.

- **`GET /v1/network-health`**  
  Returns the health status of the node.

- **`GET /v1/network-config`**  
  Returns the network configuration details.

- **`GET /v1/version`**  
  Returns the current version of the gateway

- **`GET /v1/getmessage`**  
  Generates and returns a message for the user to sign based on the provided encryption token.

