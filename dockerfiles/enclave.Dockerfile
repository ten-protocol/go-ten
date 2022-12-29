FROM ghcr.io/edgelesssys/ego-dev:latest

# on the container:
#   /home/obscuro/data       contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
#

# setup container data structure
RUN mkdir -p /home/obscuro/data && mkdir -p /home/obscuro/go-obscuro

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/obscuro/go-obscuro
COPY go.mod .
COPY go.sum .
RUN go mod download

# COPY the source code as the last step
COPY . .

# build binary
WORKDIR /home/obscuro/go-obscuro/go/enclave/main
RUN ego-go build && ego sign main

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 11000