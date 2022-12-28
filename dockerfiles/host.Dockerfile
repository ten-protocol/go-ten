# build the obscuro host image
#
FROM golang:1.17-alpine

# set the base libs to build / run
RUN apk add build-base bash
ENV CGO_ENABLED=1

# create the base directory
RUN mkdir /home/go-obscuro

# cache the go mod packaging
COPY ./go.mod /home/go-obscuro
COPY ./go.sum /home/go-obscuro
WORKDIR /home/go-obscuro
RUN go get -d -v ./...

# make sure the all code is available
COPY . /home/go-obscuro

# build the gethnetwork exec
WORKDIR /home/go-obscuro/go/host/main
RUN go build

# expose the http and the ws ports to the host
EXPOSE 8025 9000