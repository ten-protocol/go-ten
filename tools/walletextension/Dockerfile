FROM golang:1.18-alpine

RUN apk add build-base

RUN mkdir /obscuro-gateway
WORKDIR /obscuro-gateway
# todo: minimize the image size using build steps like we do on the host/enclave images
COPY . .

RUN cd tools/walletextension/main && go build -o obscuro-gateway .

EXPOSE 3000