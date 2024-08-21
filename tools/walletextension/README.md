
# Ten Gateway Documentation

For a comprehensive overview, refer to the [official documentation](https://docs.ten.xyz/docs/tools-infrastructure/hosted-gateway).

## Running the Gateway Locally

### Backend

To run the backend locally, first build it using the `go build` command. Navigate to the `tools/walletextension/main` folder and use the following commands to build for your respective operating system:

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

The binaries will be available in the `tools/walletextension/bin` directory. Run the compiled binary to start the backend.

### Frontend

Once the backend is running, navigate to the `tools/walletextension/frontend` directory and execute the following commands:

```bash
npm install
npm run dev
```

The frontend will be accessible on `http://localhost:80`.

## HTTP Endpoints

Ten Gateway exposes several HTTP endpoints for interaction:

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

