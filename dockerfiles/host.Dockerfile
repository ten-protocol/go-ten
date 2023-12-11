# Build Stages:
# system = prepares the "OS" by downloading required binaries
# get-dependencies = downloads the go modules using the prepared system
# build-host = copies over the source code and builds the binaries using a compiler cache
# final = copies over only the executables in an alpine image that doesn't have any additional load.

FROM golang:1.20-alpine3.18 as system
# set the base libs to build / run
RUN apk add build-base bash
ENV CGO_ENABLED=1

FROM system as get-dependencies
# create the base directory
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM get-dependencies as build-host
# make sure the all code is available
COPY . .

WORKDIR /home/obscuro/go-obscuro/go/host/main

# Build the host executable. Mount cross image build cache to speed up for incremental changes.
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

# Trigger another build stage to remove unnecessary files.
FROM alpine:3.18

# Copy over just the binary from the previous build stage into this one.
COPY --from=build-host \
    /home/obscuro/go-obscuro/go/host/main /home/obscuro/go-obscuro/go/host/main
    
WORKDIR /home/obscuro/go-obscuro/go/host/main

# expose the http and the ws ports to the host
EXPOSE 8025 9000