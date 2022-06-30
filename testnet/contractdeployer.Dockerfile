# deploys one contract and outputs the address
#
FROM golang:1.17-alpine

# set the base libs to build / run
RUN apk add build-base bash git
ENV CGO_ENABLED=1

# create the base directory
RUN mkdir /home/go-obscuro

# cache the go mod packaging
COPY ./go.mod /home/go-obscuro
COPY ./go.sum /home/go-obscuro
WORKDIR /home/go-obscuro
RUN go get -d -v ./...

# make sure the geth network code is available
COPY . /home/go-obscuro

# build the contract deployer exec
WORKDIR /home/go-obscuro/tools/contractdeployer/main
RUN go build
