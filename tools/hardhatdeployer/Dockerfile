# Build Stages:
# system = prepares the "OS" by downloading required binaries
# get-dependencies = downloads the go modules using the prepared system
# build-wallet = copies over the source code and builds the wallet binaries using a compiler cache
# install-npm-deps = copies over the package.json and starts downloading dependencies. 
#                    This does not depend on other stages and is executed in parallel.
# final = copies over the solidity/hardhat source, takes the installed project from the 'install-npm-deps' and the final wallet executable
#         in a lightweight image.

FROM golang:1.20-alpine3.18 as system

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

# Standalone stage to perform npm install
FROM node:16-alpine as install-npm-deps

COPY ./contracts/package.json /home/obscuro/go-obscuro/contracts/package.json
WORKDIR /home/obscuro/go-obscuro/contracts
RUN npm install

# Final stage. Compiles the solidity contracts to cache them.
FROM node:16-alpine

COPY ./contracts/ /home/obscuro/go-obscuro/contracts/
WORKDIR /home/obscuro/go-obscuro/contracts
RUN rm package-lock.json || true

COPY --from=install-npm-deps /home/obscuro/go-obscuro/contracts/ /home/obscuro/go-obscuro/contracts/

RUN npm config set update-notifier false
RUN npx hardhat compile

COPY --from=build-wallet /home/obscuro/go-obscuro/tools/walletextension/bin /home/obscuro/go-obscuro/tools/walletextension/bin

ENTRYPOINT [ "npx", "hardhat" ]