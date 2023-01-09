# deploys one contract and outputs the address
#
FROM golang:1.17-alpine

# set the base libs to build / run
RUN apk add build-base bash git
ENV CGO_ENABLED=1

# create the base directory
RUN mkdir /home/go-obscuro
RUN mkdir /home/go-obscuro/contracts

# cache the go mod packaging
COPY ./go.mod /home/go-obscuro
COPY ./go.sum /home/go-obscuro
WORKDIR /home/go-obscuro
RUN go get -d -v ./...

# make sure the geth network code is available
COPY . /home/go-obscuro

# build the contract deployer exec
WORKDIR /home/go-obscuro/tools/walletextension/main
RUN go build -o ../bin/wallet_extension_linux
WORKDIR /home/go-obscuro/tools/walletextension/bin

FROM node:lts-alpine

COPY --from=0 /home/go-obscuro/tools/walletextension/bin /home/go-obscuro/tools/walletextension/bin
COPY ./contracts /home/go-obscuro/contracts

WORKDIR /home/go-obscuro/contracts
RUN rm package-lock.json || true
RUN npm install
RUN npx hardhat compile
COPY ./contracts /home/go-obscuro/contracts
VOLUME /home/go-obscuro/contracts/deployments
ENTRYPOINT [ "npx", "hardhat" ]