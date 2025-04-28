# build a network of eth2 nodes using geth + pryms
# please check the workflows/manual-deploy-testnet-l1.yml for more info

# Build Stages:
# system = prepares the "OS" by downloading required binaries
# get-dependencies = downloads the go modules using the prepared system
# build-geth-binary = runs the build_geth_binary script once 'system' stage is finished
# build-geth-network = compiles the gethnetwork obscuro project when 'get-dependencies' stage is finished
# final = copies over the executables from the 'build-*' stages and prepares the final image.

# Build stage for downloading dependencies based on the core defined system
FROM golang:1.23.5-bookworm as get-dependencies

# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build stage for building the eth2 network runners. Will run in parallel and block on COPY if the build-geth-binary stage has not completed.
FROM get-dependencies as build-geth-network

# Create logs directory
RUN mkdir -p /home/obscuro/logs

COPY . /home/obscuro/go-obscuro

# build the eth2network exec
WORKDIR /home/obscuro/go-obscuro/integration/eth2network/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

# Create volume for logs
VOLUME /home/obscuro/logs

# expose the http and the ws ports to the host
EXPOSE 12000 12100 12200 12300 12400 12500 12600
ENTRYPOINT ["/home/obscuro/go-obscuro/integration/eth2network/main/main"]