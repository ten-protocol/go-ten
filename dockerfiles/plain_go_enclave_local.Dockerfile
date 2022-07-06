FROM golang:1.17-alpine

# on the container:
#   /data                    contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
#
# Note: ego uses a virtual file system mount to map data directory to /data inside the enclave,
#   for this non-ego build I'm using /data as the data dir to preserve /data folder in paths inside enclave
RUN mkdir /data
RUN mkdir /home/obscuro
RUN mkdir /home/obscuro/go-obscuro

# build the enclave from the current branch
COPY . /home/obscuro/go-obscuro
WORKDIR /home/obscuro/go-obscuro/go/enclave/main
RUN apk add build-base
ENV CGO_ENABLED=1
# Download all the dependencies
RUN go get -d -v ./...
# Install the package
RUN go install -v ./...
EXPOSE 11000
ENTRYPOINT ["main"]