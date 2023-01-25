# build a network of geth nodes
# please check the workflows/manual-deploy-testnet-l1.yml for more info

# Build Stages:
# system = prepares the "OS" by downloading required binaries
# get-dependencies = downloads the go modules using the prepared system
# build-geth-binary = runs the build_geth_binary script once 'system' stage is finished
# build-geth-network = compiles the gethnetwork obscuro project when 'get-dependencies' stage is finished
# final = copies over the executables from the 'build-*' stages and prepares the final image.

FROM golang:1.17-alpine as system

# set the base libs to build / run
RUN apk add build-base bash git linux-headers
ENV CGO_ENABLED=1


# Build stage for downloading dependencies based on the core defined system
FROM system as get-dependencies
# create the base directory
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build stage for building the geth binary. This runs in parallel to get-dependencies as it does not depend on it.
FROM system as build-geth-binary
# make sure the geth network code is available
COPY ./integration/gethnetwork/ /home/obscuro/go-obscuro/integration/gethnetwork/

# reset any previous geth build
WORKDIR /home/obscuro/go-obscuro/integration/gethnetwork/
RUN rm -rf /home/obscuro/go-obscuro/integration/.build/geth_bin/ && rm -rf /home/obscuro/go-obscuro/integration/gethnetwork/geth_bin/
RUN ./build_geth_binary.sh --version=v1.10.17


# Build stage for building the geth network runners. Will run in parallel and block on COPY if the build-geth-binary stage has not completed. 
FROM get-dependencies as build-geth-network

COPY . /home/obscuro/go-obscuro

# build the gethnetwork exec
WORKDIR /home/obscuro/go-obscuro/integration/gethnetwork/main
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

# Use core system as gethnetwork requires access to the golang environment.
FROM system

#Move over from the finalized build 
COPY --from=build-geth-network \
    /home/obscuro/go-obscuro/integration/gethnetwork /home/obscuro/go-obscuro/integration/gethnetwork

COPY --from=build-geth-binary \
    /home/obscuro/go-obscuro/integration/.build/geth_bin/ /home/obscuro/go-obscuro/integration/.build/geth_bin/

# expose the http and the ws ports to the host
EXPOSE 8025 8026 9000 9001
ENTRYPOINT ["/home/obscuro/go-obscuro/integration/gethnetwork/main/main", "--numNodes=2", "--startPort=8000","--websocketStartPort=9000"]
