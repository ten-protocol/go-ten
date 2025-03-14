#!/bin/bash

# Install Go 1.23.5
echo "Installing Go 1.23.5..."
wget https://go.dev/dl/go1.23.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.5.linux-amd64.tar.gz
rm go1.23.5.linux-amd64.tar.gz

# Ensure Go is in the PATH
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Verify Go installation
go version

go install github.com/ethereum/go-ethereum/cmd/abigen@v1.13.15;
nvm install 18;
cd contracts;
npm install;
cd ..
sudo apt-get update;
sudo apt-get install -y protobuf-compiler;
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest;
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest;
