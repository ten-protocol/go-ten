#!/bin/bash

# Download and install Go 1.22.11
curl -OL https://go.dev/dl/go1.22.11.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.11.linux-amd64.tar.gz
rm go1.22.11.linux-amd64.tar.gz

go install github.com/ethereum/go-ethereum/cmd/abigen@v1.13.15;
nvm install 18;
cd contracts;
npm install;
cd ..
sudo apt-get update;
sudo apt-get install -y protobuf-compiler;
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest;
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest;
