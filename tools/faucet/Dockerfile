# Build stage for downloading dependencies based on the core defined system
FROM golang:1.23-bullseye as get-dependencies

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build stage for building the eth2 network runners. Will run in parallel and block on COPY if the build-geth-binary stage has not completed.
FROM get-dependencies as build-faucet

COPY . /home/obscuro/go-obscuro

# build the faucet exec
WORKDIR /home/obscuro/go-obscuro/tools/faucet/cmd
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o faucet

EXPOSE 80