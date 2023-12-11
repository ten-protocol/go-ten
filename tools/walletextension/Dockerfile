## Build Stages:
# system = prepares the "OS" by downloading required binaries
# get-dependencies = downloads the go modules using the prepared system
# build-wallet = copies over the source code and builds the binaries using a compiler cache
# final = copies over only the executables in an alpine image that doesn't have any additional load.

FROM golang:1.20-alpine3.18 as system

# set the base libs to build / run
RUN apk add build-base bash git
ENV CGO_ENABLED=1
ARG TESTNET_TYPE

# Standard build stage that initializes the go dependencies
FROM system as get-dependencies
# create the base directory
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro


# Install Node.js and npm (needed for frontend)
RUN apk add --update nodejs npm

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build stage that will create a wallet extension executable
FROM get-dependencies as build-wallet
# make sure the geth network code is available
COPY . /home/obscuro/go-obscuro

# Create .env file for frontend
WORKDIR /home/obscuro/go-obscuro/tools/walletextension/frontend
RUN if [ "$TESTNET_TYPE" = "dev-testnet" ]; then \
      echo "NEXT_PUBLIC_API_GATEWAY_URL=https://dev-testnet.obscu.ro" > .env; \
    elif [ "$TESTNET_TYPE" = "uat-testnet" ]; then \
      echo "NEXT_PUBLIC_API_GATEWAY_URL=https://uat-testnet.obscu.ro" > .env; \
    elif [ "$TESTNET_TYPE" = "sepolia-testnet" ]; then \
      echo "NEXT_PUBLIC_API_GATEWAY_URL=https://testnet.obscu.ro" > .env; \
    else \
      echo "NEXT_PUBLIC_API_GATEWAY_URL=http://127.0.0.1:3000" > .env; \
    fi
# Run npm build for frontend
RUN npm install
RUN npm run build

# build the gateway executable
WORKDIR /home/obscuro/go-obscuro/tools/walletextension/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o ../bin/wallet_extension_linux

# Lightweight final build stage. Includes bare minimum to start wallet extension
FROM alpine:3.18

COPY --from=build-wallet /home/obscuro/go-obscuro/tools/walletextension/bin /home/obscuro/go-obscuro/tools/walletextension/bin
WORKDIR /home/obscuro/go-obscuro/tools/walletextension/bin
