# build a network of geth nodes
# please check the workflows/manual-deploy-testnet-l1.yml for more info
#
FROM golang:1.17-alpine

# set the base libs to build / run
RUN apk add build-base bash git linux-headers
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

# reset any previous geth build
WORKDIR /home/go-obscuro/integration/gethnetwork/
RUN rm -rf /home/go-obscuro/integration/.build/geth_bin/ && rm -rf /home/go-obscuro/integration/gethnetwork/geth_bin/
RUN ./build_geth_binary.sh --version=v1.10.17

# build the gethnetwork exec
WORKDIR /home/go-obscuro/integration/gethnetwork/main
RUN go build

# expose the http and the ws ports to the host
EXPOSE 8025 8026 9000 9001
ENTRYPOINT ["/home/go-obscuro/integration/gethnetwork/main/main", "--numNodes=2", "--startPort=8000","--websocketStartPort=9000"]
