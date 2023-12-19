# Build stages:
# system = prepares the "OS" by downloading required binaries, installs dlv.
# get-dependencies = using the "system" downloads the go modules.
# build-enclave = copies over the source and builds the enclave using a go compiler cache
# final = using the base system copies over only the enclave executable and creates the final image without source and dependencies. 

FROM golang:1.20-alpine3.18 as system

# install build utils
RUN apk add build-base
ENV CGO_ENABLED=1
RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.2

FROM system as get-dependencies
# setup container data structure
RUN mkdir -p /enclavedata && mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM get-dependencies as build-enclave
# COPY the source code as the last step
COPY . .

# build the enclave from the current branch
WORKDIR /home/obscuro/go-obscuro/go/enclave/main
# Install the package
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

# Debug image which installs and runs delve and runs the enclave package without eGo
#
# Final container folder structure:
#   /data                    contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
#
# Note: ego uses a virtual file system mount to map data directory to /data inside the enclave,
#   for this non-ego build I'm using /data as the data dir to preserve /data folder in paths inside enclave
#

EXPOSE 11000