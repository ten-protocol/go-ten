# 👛 The Obscuro wallet extension

See the documentation [here](https://docsstage.obscu.ro/testnet/wallet-extension.html).

# Developer notes

## Compiling the binary for your own computer

The wallet extension binary can be built for your own computer by running:

```
go build -o ../bin/wallet_extension .
```

The binary will be created in the `tools/walletextension/bin` folder.

## Compiling the binary for other common architectures

The precompiled binaries for macOS ARM64, macOS AMD64, Windows AMD64 and Linux AMD64 can be built by running the 
following commands from the `tools/walletextension/main` folder.

```
env GOOS=darwin GOARCH=amd64 go build -o ../bin/wallet_extension_macos_amd64 .
env GOOS=darwin GOARCH=arm64 go build -o ../bin/wallet_extension_macos_arm64 .
env GOOS=windows GOARCH=amd64 go build -o ../bin/wallet_extension_win_amd64.exe .
env GOOS=linux GOARCH=amd64 go build -o ../bin/wallet_extension_linux_amd64 .
```

The binaries will be created in the `tools/walletextension/bin` folder.
