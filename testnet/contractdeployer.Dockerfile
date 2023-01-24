# runs obscuro scan
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

FROM get-dependencies as build-deployer
# make sure the geth network code is available
COPY . /home/obscuro/go-obscuro

# build the contract deployer exec
WORKDIR /home/obscuro/go-obscuro/tools/contractdeployer/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

FROM alpine:3.17

COPY --from=build-deployer\
    /home/obscuro/go-obscuro/tools/contractdeployer/main /home/obscuro/go-obscuro/tools/contractdeployer/main
    
WORKDIR /home/obscuro/go-obscuro

