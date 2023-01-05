FROM golang:1.17-alpine

# Debug image which installs and runs delve and runs the enclave package without eGo
#
# on the container:
#   /data                    contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
#
# Note: ego uses a virtual file system mount to map data directory to /data inside the enclave,
#   for this non-ego build I'm using /data as the data dir to preserve /data folder in paths inside enclave
#

# install build utils
RUN apk add build-base
ENV CGO_ENABLED=1
RUN go install github.com/go-delve/delve/cmd/dlv@v1.9.1

# setup container data structure
RUN mkdir -p /data && mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# COPY the source code as the last step
COPY . .

# build the enclave from the current branch
WORKDIR /home/obscuro/go-obscuro/go/enclave/main
# Install the package
RUN go get -v ./...

EXPOSE 11000