# runs ObscuroScan
#
FROM ghcr.io/edgelesssys/ego-dev:latest

# set the base libs to build / run
RUN apt-get -y install software-properties-common
RUN yes | ego install az-dcap-client
ENV CGO_ENABLED=1

# create the base directory
RUN mkdir /home/go-obscuro

# cache the go mod packaging
COPY ./go.mod /home/go-obscuro
COPY ./go.sum /home/go-obscuro
WORKDIR /home/go-obscuro
RUN go get -d -v ./...

# make sure the ObscuroScan code is available
COPY . /home/go-obscuro

# build the ObscuroScan exec
WORKDIR /home/go-obscuro/tools/obscuroscan/main
RUN ego-go build && ego sign main

WORKDIR /home/go-obscuro
