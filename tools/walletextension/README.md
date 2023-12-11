# The Ten gateway

See the documentation [here](https://docs.obscu.ro/wallet-extension/wallet-extension/).

## Developer notes

Running gateway frontend locally requires building static files first.
To do that, run `npm run build` in `tools/walletextension/frontend` folder.

The precompiled binaries for macOS ARM64, macOS AMD64, Windows AMD64 and Linux AMD64 can be built by running the 
following commands from the `tools/walletextension/main` folder:

```
env GOOS=darwin GOARCH=amd64 go build -o ../bin/wallet_extension_macos_amd64 .
env GOOS=darwin GOARCH=arm64 go build -o ../bin/wallet_extension_macos_arm64 .
env GOOS=windows GOARCH=amd64 go build -o ../bin/wallet_extension_win_amd64.exe .
env GOOS=linux GOARCH=amd64 go build -o ../bin/wallet_extension_linux_amd64 .
```

The binaries will be created in the `tools/walletextension/bin` folder.

### Structure

This package follows the same structure of `host` and `enclave`.

It uses a container to wrap the services that are required to allow the wallet extension to fulfill the business logic.

### Running Wallet Extension with Docker

To build a docker image use docker build command. Please note that you need to run it from the root of the repository.
To run the container you can use `./docker_run.sh`. You can add parameters to the script, and they are passed to the wallet extension 
(example: `-host=0.0.0.0` to be able to access wallet extension endpoints via localhost).


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
