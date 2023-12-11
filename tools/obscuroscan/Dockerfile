# Build Stages:
# system = prepares the "OS" by downloading required binaries
# get-dependencies = downloads the go modules using the prepared system
# build-obscuroscan = copies over the source code and builds the binaries using a compiler cache
# final = copies over only the executables in an alpine image that doesn't have any additional load.

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

FROM get-dependencies as build-obscuroscan
# make sure the geth network code is available
COPY . /home/obscuro/go-obscuro

# build the contract deployer exec
WORKDIR /home/obscuro/go-obscuro/tools/obscuroscan/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

FROM alpine:3.18

COPY --from=build-obscuroscan /home/obscuro/go-obscuro/tools/obscuroscan/main /home/obscuro/go-obscuro/tools/obscuroscan/main
WORKDIR /home/obscuro/go-obscuro
