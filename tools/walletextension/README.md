# 👛 The Obscuro wallet extension

See the documentation [here](https://docs.obscu.ro/wallet-extension/wallet-extension/).

## Developer notes

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

To build a docker image use `./docker_build.sh` script and for running it locally
you can use `./docker_run.sh`. You can add parameters to the script, and they are passed to wallet extension 
(example: `-host=0.0.0.0` to be able to access wallet extension endpoints via localhost).
