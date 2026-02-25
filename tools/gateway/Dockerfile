## Build Stages:
# system = prepares the "OS" by downloading required binaries
# get-dependencies = downloads the go modules using the prepared system
# build-wallet = copies over the source code and builds the binaries using a compiler cache
# final = copies over only the executables in an alpine image that doesn't have any additional load.

FROM golang:1.25-alpine AS system

# set the base libs to build / run
RUN apk add build-base bash git
ENV CGO_ENABLED=1
ARG TESTNET_TYPE

# Standard build stage that initializes the go dependencies
FROM system AS get-dependencies
# create the base directory
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro


# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build stage that will create a gateway executable
FROM get-dependencies AS build-wallet
# make sure the geth network code is available
COPY . /home/obscuro/go-obscuro

# build the gateway executable
WORKDIR /home/obscuro/go-obscuro/tools/gateway/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o ../bin/gateway_linux

# Lightweight final build stage. Includes bare minimum to start gateway
FROM alpine:3.18


# Set a writable working directory
WORKDIR /home/obscuro
RUN mkdir -p /home/obscuro && chmod -R 777 /home/obscuro

# copy over the gateway executable
COPY --from=build-wallet /home/obscuro/go-obscuro/tools/gateway/bin /home/obscuro/go-obscuro/tools/gateway/bin

# copy over the .sql migration files
COPY --from=build-wallet /home/obscuro/go-obscuro/tools/gateway/storage/database /home/obscuro/go-obscuro/tools/gateway/storage/database

# copy over the entrypoint script
COPY --from=build-wallet /home/obscuro/go-obscuro/tools/gateway/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

