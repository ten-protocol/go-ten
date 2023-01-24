# deploys one contract and outputs the address
#
FROM golang:1.17-alpine as system

# set the base libs to build / run
RUN apk add build-base bash git
ENV CGO_ENABLED=1

# Standard build stage that initializes the go dependencies
FROM system as get-dependencies
# create the base directory
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build stage that will create a wallet extension executable
FROM get-dependencies as build-wallet
# make sure the geth network code is available
COPY . /home/obscuro/go-obscuro

# build the contract deployer exec
WORKDIR /home/obscuro/go-obscuro/tools/walletextension/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o ../bin/wallet_extension_linux

# Lightweight final build stage. Includes bare minimum to start wallet extension
FROM alpine:3.17

COPY --from=build-wallet /home/obscuro/go-obscuro/tools/walletextension/bin /home/obscuro/go-obscuro/tools/walletextension/bin
WORKDIR /home/obscuro/go-obscuro/tools/walletextension/bin
ENTRYPOINT [ "./wallet_extension_linux" ]
