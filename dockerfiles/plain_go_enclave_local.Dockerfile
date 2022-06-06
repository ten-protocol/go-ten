FROM golang:1.17-alpine

# build the enclave from the current branch
RUN mkdir /home/obscuro-playground
COPY . /home/obscuro-playground
WORKDIR /home/obscuro-playground/go/obscuronode/enclave/main
RUN apk add build-base
ENV CGO_ENABLED=1
# Download all the dependencies
RUN go get -d -v ./...
# Install the package
RUN go install -v ./...
EXPOSE 11000
ENTRYPOINT ["main"]