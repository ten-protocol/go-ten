FROM golang:1.17-alpine

# on the container:
#   /data                    contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
#
# Note: on ego it maps /home/obscuro/data to /data inside the enclave,
#   using /data here should allow us to preserve /data folder in paths inside enclave)
RUN mkdir /data
RUN mkdir /home/obscuro
RUN mkdir /home/obscuro/go-obscuro

# build the enclave from the current branch
COPY . /home/obscuro/go-obscuro
WORKDIR /home/obscuro/go-obscuro/go/obscuronode/enclave/main
RUN apk add build-base
ENV CGO_ENABLED=1
# Download all the dependencies
RUN go get -d -v ./...
# Install the package
RUN go install -v ./...
EXPOSE 11000
ENTRYPOINT ["main"]