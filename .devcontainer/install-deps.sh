#!/bin/bash

# Skip Go installation as it's already provided by the devcontainer image
echo "Using existing Go installation..."
go version

# Ensure Go is in the PATH (should already be set)
export PATH=$PATH:/usr/local/go/bin

# Set CGO environment variables to fix cross-compilation issues
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

# Install Go tools with proper CGO configuration
echo "Installing Go tools..."
CGO_ENABLED=0 go install github.com/ethereum/go-ethereum/cmd/abigen@v1.13.15

# Install Node.js 20
echo "Installing Node.js 20..."
nvm install 20

# Install npm dependencies
echo "Installing npm dependencies..."
cd contracts
npm install
cd ..

# Update package list and install protobuf compiler
echo "Installing system packages..."
sudo apt-get update
sudo apt-get install -y protobuf-compiler

# Install protobuf Go plugins
echo "Installing protobuf Go plugins..."
CGO_ENABLED=0 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
CGO_ENABLED=0 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install abigen again (latest version)
CGO_ENABLED=0 go install github.com/ethereum/go-ethereum/cmd/abigen@latest

echo "Installation completed successfully!"

